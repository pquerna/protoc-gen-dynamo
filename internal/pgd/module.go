package pgd

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/davecgh/go-spew/spew"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
	dynamopb "github.com/pquerna/protoc-gen-dynamo/dynamo"
)

const (
	moduleName    = "dynamo"
	version       = "0.1.0"
	commentFormat = `// Code generated by protoc-gen-%s v%s. DO NOT EDIT.
// source: %s
`
)

type Module struct {
	*pgs.ModuleBase
	ctx pgsgo.Context
}

var _ pgs.Module = (*Module)(nil)

func New() pgs.Module {
	return &Module{ModuleBase: &pgs.ModuleBase{}}
}

func (m *Module) InitContext(ctx pgs.BuildContext) {
	m.ModuleBase.InitContext(ctx)
	m.ctx = pgsgo.InitContext(ctx.Parameters())
}

func (m *Module) Name() string {
	return moduleName
}

func (m *Module) Execute(targets map[string]pgs.File, pkgs map[string]pgs.Package) []pgs.Artifact {
	for _, f := range targets {
		msgs := f.AllMessages()
		if n := len(msgs); n == 0 {
			m.Debugf("No messagess in %v, skipping", f.Name())
			continue
		}
		m.processFile(f)
	}
	return m.Artifacts()
}

func (m *Module) processFile(f pgs.File) {
	out := bytes.Buffer{}
	err := m.applyTemplate(&out, f)
	if err != nil {
		m.Logf("couldn't apply template: %s", err)
		m.Fail("code generation failed")
	} else {
		generatedFileName := m.ctx.OutputPath(f).SetExt(fmt.Sprintf(".%s.go", moduleName)).String()
		m.AddGeneratorFile(generatedFileName, out.String())
	}
}

const (
	dynamoPkg  = "github.com/aws/aws-sdk-go/service/dynamodb"
	protoPkg   = "github.com/golang/protobuf/proto"
	awsPkg     = "github.com/aws/aws-sdk-go/aws"
	strconvPkg = "strconv"
	stringsPkg = "strings"
	fmtPkg     = "fmt"
	timePkg    = "time"

	timestampType = "google.protobuf.Timestamp"
)

func (m *Module) applyTemplate(buf *bytes.Buffer, in pgs.File) error {
	pkgName := m.ctx.PackageName(in).String()
	importPath := m.ctx.ImportPath(in).String()
	protoFileName := in.Name().String()

	f := jen.NewFilePathName(importPath, pkgName)
	f.HeaderComment(fmt.Sprintf(commentFormat, moduleName, version, protoFileName))

	f.ImportName(dynamoPkg, "dynamodb")
	f.ImportName(awsPkg, "aws")
	f.ImportName(protoPkg, "proto")
	f.ImportName(strconvPkg, "strconv")
	f.ImportName(fmtPkg, "fmt")
	f.ImportName(stringsPkg, "strings")
	f.ImportName(timePkg, "time")

	// https://godoc.org/github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute#Marshaler
	// https://godoc.org/github.com/guregu/dynamo#Marshaler
	err := m.applyMarshal(f, in)
	if err != nil {
		return err
	}

	// https://godoc.org/github.com/guregu/dynamo#Unmarshaler
	// UnmarshalDynamo(av *dynamodb.AttributeValue) error
	err = m.applyUnmarshal(f, in)
	if err != nil {
		return err
	}

	err = m.applyKeyFuncs(f, in)
	if err != nil {
		return err
	}

	return f.Render(buf)
}

type avType string

const (
	avt_bytes      avType = "B"
	avt_bool       avType = "BOOL"
	avt_byte_set   avType = "BS"
	avt_list       avType = "L"
	avt_map        avType = "M"
	avt_number     avType = "N"
	avt_number_set avType = "NS"
	avt_null       avType = "NULL"
	avt_string     avType = "S"
	avt_string_set avType = "SS"
)

func getAVType(field pgs.Field, fext *dynamopb.DynamoFieldOptions) avType {
	isArray := field.Type().ProtoLabel() == pgs.Repeated
	pt := field.Type().ProtoType()

	if isArray {
		if !fext.Type.Set {
			return avt_list
		}
		switch {
		case pt.IsInt() || pt == pgs.DoubleT || pt == pgs.FloatT:
			return avt_number_set
		case pt == pgs.StringT:
			return avt_string_set
		case pt == pgs.BytesT:
			return avt_byte_set
		case pt == pgs.EnumT:
			return avt_number_set
		}
	} else {
		switch {
		case pt.IsInt() || pt == pgs.DoubleT || pt == pgs.FloatT:
			return avt_number
		case pt == pgs.BoolT:
			return avt_bool
		case pt == pgs.StringT:
			return avt_string
		case pt == pgs.BytesT:
			return avt_bytes
		case pt == pgs.MessageT:
			return avt_map
		case pt == pgs.EnumT:
			return avt_number
		}
	}
	panic(fmt.Sprintf("getAVType: failed to determine dynamodb type: %T %+v", field, fext.Type))
}

func fieldByName(msg pgs.Message, name string) pgs.Field {
	for _, f := range msg.Fields() {
		if f.Name().LowerSnakeCase().String() == name {
			return f
		}
	}
	panic(fmt.Sprintf("Failed to find field %s on %s", name, msg.FullyQualifiedName()))
}

type namedKey struct {
	name string
	constant string
	fields []string
}

func getVersionField(msg pgs.Message) pgs.Field {
	fn := "updated_at"
	for _, f := range msg.Fields() {
		if f.Name().LowerSnakeCase().String() == fn {
			return f
		}
	}
	return nil
}

func (m *Module) applyVersionFuncs(msg pgs.Message, f *jen.File) error {
	structName := m.ctx.Name(msg)

	field := getVersionField(msg)
	if field == nil {
		// No version field, don't apply version func
		return nil
	}

	fn := field.Name().String()
	srcName := field.Name().UpperCamelCase().String()

	var stmts []jen.Code
	d := field.Descriptor().TypeName
	if d == nil {
		return errors.New(fmt.Sprintf("Failed to find field descriptor for %s on %s", fn, msg.FullyQualifiedName()))
	}
	if !strings.HasSuffix(*d, timestampType) {
		return errors.New(fmt.Sprintf("Field descriptor for %s on %s is not a timestamp", fn, msg.FullyQualifiedName()))
	}
	// 	err := p.UpdatedAt.CheckValid()
	//	if err != nil {
	//		return 0, err
	//	}
	//	t := p.UpdatedAt.AsTime()
	// return t.UnixNano(), nil
	stmts = append(stmts, jen.List(jen.Err()).Op(":=").Id("p").Dot(srcName).Dot("CheckValid").Call())
	stmts = append(stmts, jen.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.List(jen.Lit(0), jen.Err()))))
	stmts = append(stmts, jen.List(jen.Id("t")).Op(":=").Id("p").Dot(srcName).Dot("AsTime").Call())
	stmts = append(stmts, jen.Return(jen.List(jen.Id("t").Dot("UnixNano").Call(), jen.Nil())))

	if len(stmts) == 0 {
		return errors.New("version: numeric or timestamp type is required")
	}

	f.Func().Params(
		jen.Id("p").Op("*").Id(structName.String()),
	).Id("Version").Params().Parens(jen.List(jen.Int64(), jen.Error())).Block(stmts...).Line()

	return nil
}

func (m *Module) applyKeyFuncs(f *jen.File, in pgs.File) error {
	const stringBuffer = "sb"
	for _, msg := range in.AllMessages() {
		structName := m.ctx.Name(msg)
		mext := dynamopb.DynamoMessageOptions{}
		ok, err := msg.Extension(dynamopb.E_Msg, &mext)
		if err != nil {
			m.Logf("Parsing dynamo.msg.disabled failed: %s", err)
			m.Fail("code generation failed")
		}
		if ok && mext.Disabled {
			m.Logf("dynamo.msg disabled for %s", structName)
			continue
		}

		keys := []namedKey{ }
		// pk, sk, gsi1pk, gsi1sk
		for _, ck := range mext.Key {

			pkName := "PartitionKey"
			if ck.GlobalSecondaryIndex != 0 {
				pkName = fmt.Sprintf("Gsi%dPkKey", ck.GlobalSecondaryIndex)
			}

			skName := "SortKey"
			if ck.GlobalSecondaryIndex != 0 {
				skName = fmt.Sprintf("Gsi%dSkKey", ck.GlobalSecondaryIndex)
			}

			keys = append(keys,
				namedKey{
					name: pkName,
					fields: ck.PkFields,
				})

			keys = append(keys,
				namedKey{
					name: skName,
					constant: ck.SkConst,
					fields: ck.SkFields,
				})
		}

		if err := m.applyVersionFuncs(msg, f); err != nil {
			m.Logf("Generating version funcs failed: %s", err)
			m.Fail("code generation failed")
		}

		for _, key := range keys {
			stmts := []jen.Code{}
			if key.constant != "" {
				stmts = append(stmts,
					jen.Return(jen.Lit(key.constant)),
				)
			} else {
				stmts = append(stmts,
					jen.Op("var").Id("sb").Qual(stringsPkg, "Builder"),
				)
				stmts = generateKeyStringer(msg, stmts, key.fields, stringBuffer)
				stmts = append(stmts,
					jen.Return(jen.Id(stringBuffer).Dot("String").Call()),
				)
			}

			f.Func().Params(
				jen.Id("p").Op("*").Id(structName.String()),
			).Id(key.name).Params().List(jen.String()).Block(
				stmts...,
			).Line()

			params := []jen.Code{}
			d := jen.Dict{}
			for _, fn := range key.fields {
				field := fieldByName(msg, fn)
				typ := m.ctx.Type(field)
				params = append(params, jen.Id(field.Name().LowerCamelCase().String()).Id(typ.String()))
				d[jen.Id(field.Name().UpperCamelCase().String())] = jen.Id(field.Name().LowerCamelCase().String())
			}

			f.Func().Id(structName.String() + key.name).Params(params...).List(jen.String()).Block(
				jen.Return(jen.Call(jen.Op("&").Id(structName.String()).Values(d)).Dot(key.name).Call()),
			).Line()
		}
	}
	return nil
}

func generateKeyStringer(msg pgs.Message, stmts []jen.Code, fields []string, stringBuffer string) []jen.Code {
	stmts = append(stmts, jen.Id(stringBuffer).Dot("Reset").Call())

	sep := ":"
	prefix := msg.Name().LowerSnakeCase().String()
	stmts = append(stmts, jen.List(jen.Id("_"), jen.Id("_")).Op("=").Id(stringBuffer).Dot("WriteString").Call(
		jen.Lit(prefix+sep),
	))

	first := true
	for _, fn := range fields {
		field := fieldByName(msg, fn)
		pt := field.Type().ProtoType()
		srcName := field.Name().UpperCamelCase().String()
		if !first {
			stmts = append(stmts, jen.List(jen.Id("_"), jen.Id("_")).Op("=").Id(stringBuffer).Dot("WriteString").Call(
				jen.Lit(sep),
			))
		}
		first = false
		switch {
		case pt == pgs.StringT:
			stmts = append(stmts, jen.List(jen.Id("_"), jen.Id("_")).Op("=").Id(stringBuffer).Dot("WriteString").Call(
				jen.Id("p").Dot(srcName),
			))
		case pt.IsNumeric() || pt == pgs.EnumT:
			fmtCall := numberFormatStatement(pt, jen.Id("p").Dot(srcName))
			stmts = append(stmts, jen.List(jen.Id("_"), jen.Id("_")).Op("=").Id(stringBuffer).Dot("WriteString").Call(
				fmtCall,
			))
		default:
			panic(fmt.Sprintf("Compound key: unsupported type: %s", pt.String()))
		}
	}
	return stmts
}

const (
	valueField   = "value"
	deletedField = "deleted"
	typeField    = "typ"
)

func (m *Module) applyMarshal(f *jen.File, in pgs.File) error {
	for _, msg := range in.AllMessages() {
		structName := m.ctx.Name(msg)
		mext := dynamopb.DynamoMessageOptions{}
		ok, err := msg.Extension(dynamopb.E_Msg, &mext)
		if err != nil {
			m.Logf("Parsing dynamo.msg.disabled failed: %s", err)
			m.Fail("code generation failed")
		}
		if ok && mext.Disabled {
			m.Logf("dynamo.msg disabled for %s", structName)
			continue
		}

		// https://godoc.org/github.com/guregu/dynamo#Marshaler:
		// MarshalDynamo() (*dynamodb.AttributeValue, error)
		f.Func().Params(
			jen.Id("p").Op("*").Id(structName.String()),
		).Id("MarshalDynamo").Params().List(jen.Params(jen.Op("*").Qual(dynamoPkg, "AttributeValue"), jen.Id("error"))).Block(
			jen.Id("av").Op(":=").Op("&").Qual(dynamoPkg, "AttributeValue").Values(),
			jen.Id("err").Op(":=").Id("p").Dot("MarshalDynamoDBAttributeValue").Call(jen.Id("av")),
			jen.If(jen.Id("err").Op("!=").Nil()).Block(
				jen.Return(jen.Nil(), jen.Id("err")),
			),
			jen.Return(jen.Id("av"), jen.Nil()),
		).Line()

		f.Func().Params(
			jen.Id("p").Op("*").Id(structName.String()),
		).Id("MarshalDynamoItem").Params().List(jen.Params(jen.Map(jen.String()).Op("*").Qual(dynamoPkg, "AttributeValue"), jen.Id("error"))).Block(
			jen.Id("av").Op(":=").Op("&").Qual(dynamoPkg, "AttributeValue").Values(),
			jen.Id("err").Op(":=").Id("p").Dot("MarshalDynamoDBAttributeValue").Call(jen.Id("av")),
			jen.If(jen.Id("err").Op("!=").Nil()).Block(
				jen.Return(jen.Nil(), jen.Id("err")),
			),
			jen.Return(jen.Id("av").Dot("M"), jen.Nil()),
		).Line()

		stmts := []jen.Code{}
		refId := 0
		d := jen.Dict{}
		needErr := false
		needNullBoolTrue := false
		needProtoBuffer := false
		needStringBuilder := false
		const protoBuffer = "pbuf"
		const stringBuffer = "sb"
		computedKeys := make([]*dynamopb.Key, 0)
		if mext.Key != nil {
			computedKeys = append(computedKeys, mext.Key...)
		}

		if false {
			m.Log(spew.Sprint(computedKeys))
		}

		for _, ck := range computedKeys {
			needStringBuilder = true

			pkName := "pk"
			if ck.GlobalSecondaryIndex != 0 {
				pkName = fmt.Sprintf("gsi%dpk", ck.GlobalSecondaryIndex)
			}
			skName := "sk"
			if ck.GlobalSecondaryIndex != 0 {
				skName = fmt.Sprintf("gsi%dsk", ck.GlobalSecondaryIndex)
			}

			refId++
			vname := fmt.Sprintf("v%d", refId)
			stmts = append(stmts, jen.Id(vname).Op(":=").Op("&").Qual(dynamoPkg, "AttributeValue").Values())
			stmts = generateKeyStringer(msg, stmts, ck.PkFields, stringBuffer)
			stmts = append(stmts, jen.Id(vname).Dot("S").Op("=").Qual(awsPkg, "String").Call(jen.Id(stringBuffer).Dot("String").Call()))
			d[jen.Lit(pkName)] = jen.Id(vname)
			if ck.SkConst != "" {
				refId++
				vname = fmt.Sprintf("v%d", refId)
				stmts = append(stmts, jen.Id(vname).Op(":=").Op("&").Qual(dynamoPkg, "AttributeValue").Values(jen.Dict{
					jen.Id("S"): jen.Qual(awsPkg, "String").Call(jen.Lit(ck.SkConst)),
				}))
				d[jen.Lit(skName)] = jen.Id(vname)
			} else {
				refId++
				vname = fmt.Sprintf("v%d", refId)
				stmts = append(stmts, jen.Id(vname).Op(":=").Op("&").Qual(dynamoPkg, "AttributeValue").Values())
				stmts = generateKeyStringer(msg, stmts, ck.SkFields, stringBuffer)
				stmts = append(stmts, jen.Id(vname).Dot("S").Op("=").Qual(awsPkg, "String").Call(jen.Id(stringBuffer).Dot("String").Call()))
				d[jen.Lit(skName)] = jen.Id(vname)
			}
		}

		if getVersionField(msg) != nil {
			refId++
			vname := fmt.Sprintf("v%d", refId)
			// Version() (int64, error)
			stmts = append(stmts, jen.List(jen.Id(vname), jen.Id("err")).Op(":=").Id("p").Dot("Version").Call())
			stmts = append(stmts,
				jen.If(jen.Id("err").Op("!=").Nil()).Block(
					jen.Return(jen.Id("err")),
				),
			)
			d[jen.Lit("version")] = jen.Op("&").Qual(dynamoPkg, "AttributeValue").Values(jen.Dict{
				jen.Id("N"): jen.Qual(awsPkg, "String").Call(jen.Qual(strconvPkg, "FormatInt").Call(jen.Id(vname), jen.Lit(10))),
			})
		}

		typeName := fmt.Sprintf("%s.%s", msg.Package().ProtoName().String(), msg.Name())

		needProtoBuffer = true
		needErr = true
		refId++
		vname := fmt.Sprintf("v%d", refId)
		stmts = append(stmts, jen.Id(vname).Op(":=").Op("&").Qual(dynamoPkg, "AttributeValue").Values())
		stmts = append(stmts, jen.Id(protoBuffer).Dot("Reset").Call())
		stmts = append(stmts, jen.Id("err").Op("=").Id(protoBuffer).Dot("Marshal").Call(jen.Id("p")))
		stmts = append(stmts,
			jen.If(jen.Id("err").Op("!=").Nil()).Block(
				jen.Return(jen.Id("err")),
			),
		)
		stmts = append(stmts, jen.Id(vname).Dot("B").Op("=").Id(protoBuffer).Dot("Bytes").Call())
		d[jen.Lit(valueField)] = jen.Id(vname)

		refId++
		vname = fmt.Sprintf("v%d", refId)
		stmts = append(stmts, jen.Id(vname).Op(":=").Op("&").Qual(dynamoPkg, "AttributeValue").Values())
		stmts = append(stmts, jen.Id(vname).Dot("S").Op("=").Qual(awsPkg, "String").Call(jen.Lit(typeName)))
		d[jen.Lit(typeField)] = jen.Id(vname)

		for _, field := range msg.Fields() {
			fieldDescriptorName := field.Descriptor().GetTypeName()
			if strings.HasSuffix(fieldDescriptorName, timestampType) &&
				field.Name().LowerSnakeCase().String() == "deleted_at" {
				srcName := field.Name().UpperCamelCase().String()
				refId++
				vname = fmt.Sprintf("v%d", refId)
				stmts = append(stmts, jen.Id(vname).Op(":=").Op("&").Qual(dynamoPkg, "AttributeValue").Values())
				stmts = append(stmts, jen.Id(vname).Dot("BOOL").Op("=").Qual(awsPkg, "Bool").Call(jen.Id("p").Dot(srcName).Dot("IsValid").Call()))
				d[jen.Lit(deletedField)] = jen.Id(vname)
			}
		}

		for _, field := range msg.Fields() {
			fext := dynamopb.DynamoFieldOptions{}
			ok, err := field.Extension(dynamopb.E_Field, &fext)
			if err != nil {
				m.Failf("Error: Parsing dynamo.field failed for '%s': %s", field.FullyQualifiedName(), err)
			}

			if !ok {
				m.Debugf("dynamo.field.expose: skipped %s (no extension)", field.FullyQualifiedName())
				continue
			}
			if !fext.Expose {
				m.Debugf("dynamo.field.expose: skipped %s (not exposed)", field.FullyQualifiedName())
				continue
			}

			if fext.Type == nil {
				fext.Type = &dynamopb.Types{}
			}
			pt := field.Type().ProtoType()

			srcName := field.Name().UpperCamelCase().String()
			refId++
			vname := fmt.Sprintf("v%d", refId)
			arrix := fmt.Sprintf("ix%d", refId)
			arrname := fmt.Sprintf("arr%d", refId)

			isArray := field.Type().ProtoLabel() == pgs.Repeated
			if fext.Type.Set && !isArray {
				m.Failf("Error: dynamo.field.set=true, but field is not repeated / array type: '%s'.IsRepeated=%v",
					field.FullyQualifiedName(), field.Type().IsRepeated())
			}

			avt := getAVType(field, &fext)
			fieldName := jen.Lit(field.Name().LowerSnakeCase().String())
			if fext.Name != "" {
				fieldName = jen.Lit(fext.Name)
			}
			switch avt {
			case avt_bytes:
				needNullBoolTrue = true
				stmts = append(stmts, jen.Id(vname).Op(":=").Op("&").Qual(dynamoPkg, "AttributeValue").Values())
				stmts = append(stmts,
					jen.If(jen.Len(jen.Id("p").Dot(field.Name().UpperCamelCase().String())).Op("!=").Lit(0)).Block(
						jen.Id(vname).Dot("B").Op("=").Id("p").Dot(srcName),
					).Else().Block(
						jen.Id(vname).Dot("NULL").Op("=").Op("&").Id("nullBoolTrue"),
					),
				)
				d[fieldName] = jen.Id(vname)
			case avt_bool:
				d[fieldName] = jen.Op("&").Qual(dynamoPkg, "AttributeValue").Values(jen.Dict{
					jen.Id("BOOL"): jen.Op("&").Id("p").Dot(srcName),
				})
			case avt_list:
				stmts = append(stmts,
					jen.Id(arrname).Op(":=").Make(
						jen.Op("[]*").Qual(dynamoPkg, "AttributeValue"),
						jen.Lit(0),
						jen.Len(jen.Id("p").Dot(srcName)),
					),
				)

				switch {
				case pt.IsInt() || pt == pgs.DoubleT || pt == pgs.FloatT:
					fmtCall := numberFormatStatement(pt, jen.Id(arrix))
					stmts = append(stmts,
						jen.For(jen.List(jen.Id("_"), jen.Id(arrix)).Op(":=").Range().Id("p").Dot(srcName)).Block(
							jen.Id(arrname).Op("=").Append(
								jen.Id(arrname),
								jen.Op("&").Qual(dynamoPkg, "AttributeValue").Values(jen.Dict{
									jen.Id("N"): jen.Qual(awsPkg, "String").Call(
										fmtCall,
									),
								}),
							),
						),
					)
				case pt == pgs.StringT:
					stmts = append(stmts,
						jen.For(jen.List(jen.Id("_"), jen.Id(arrix)).Op(":=").Range().Id("p").Dot(srcName)).Block(
							jen.Id(arrname).Op("=").Append(
								jen.Id(arrname),
								jen.Op("&").Qual(dynamoPkg, "AttributeValue").Values(jen.Dict{
									jen.Id("S"): jen.Qual(awsPkg, "String").Call(
										jen.Id(arrix),
									),
								}),
							),
						),
					)
				default:
					m.Failf("Error: dynamo.field '%s' is repeated, but the '%s' type is not supported", field.FullyQualifiedName(), pt.String())
				}
				d[fieldName] = jen.Op("&").Qual(dynamoPkg, "AttributeValue").Values(jen.Dict{
					jen.Id("L"): jen.Id(arrname),
				})
			case avt_map:
				// if its a timestamp, keep going
				fieldDescriptorName := field.Descriptor().TypeName
				if fieldDescriptorName == nil || (fieldDescriptorName != nil && !strings.HasSuffix(*fieldDescriptorName, timestampType)) {
					m.Failf("dynamo.field: not done: avt_map type: %s / %s", field.FullyQualifiedName(), *field.Descriptor().TypeName)
					panic("applyMarshal not done: avt_map for non-timestamps")
				}
				switch {
				case fext.Type.UnixMilli, fext.Type.UnixNano, fext.Type.UnixSecond:
					var access *jen.Statement
					switch {
					case fext.Type.UnixMilli:
						// .Round(time.Millisecond).UnixNano() / time.Millisecond
						access = jen.Id("p").Dot(srcName).Dot("AsTime").Call().
							Dot("Round").Call(jen.Qual(timePkg, "Millisecond")).
							Dot("UnixNano").Call().Op("/").Int64().Call(jen.Qual(timePkg, "Millisecond"))
					case fext.Type.UnixNano:
						access = jen.Id("p").Dot(srcName).Dot("AsTime").Call().Dot("UnixNano").Call()
					case fext.Type.UnixSecond:
						access = jen.Id("p").Dot(srcName).Dot("AsTime").Call().
							Dot("Round").Call(jen.Qual(timePkg, "Second")).
							Dot("Unix").Call()
					}
					fmtCall := numberFormatStatement(pgs.Int64T, access)
					d[fieldName] = jen.Op("&").Qual(dynamoPkg, "AttributeValue").Values(jen.Dict{
						jen.Id("N"): jen.Qual(awsPkg, "String").Call(fmtCall),
					})
				default:
					m.Failf("dynamo.field: not done: applyMarshal not done: timestamps must specify the conversion type %s", field.FullyQualifiedName())
					panic("applyMarshal not done: timestamps must specify the conversion type")
				}
			case avt_number:
				fmtCall := numberFormatStatement(pt, jen.Id("p").Dot(srcName))
				d[fieldName] = jen.Op("&").Qual(dynamoPkg, "AttributeValue").Values(jen.Dict{
					jen.Id("N"): jen.Qual(awsPkg, "String").Call(fmtCall),
				})
			case avt_null:
				// avt_null: unused
			case avt_string:
				needNullBoolTrue = true
				stmts = append(stmts, jen.Id(vname).Op(":=").Op("&").Qual(dynamoPkg, "AttributeValue").Values())
				stmts = append(stmts,
					jen.If(jen.Len(jen.Id("p").Dot(field.Name().UpperCamelCase().String())).Op("!=").Lit(0)).Block(
						jen.Id(vname).Dot("S").Op("=").Qual(awsPkg, "String").Call(jen.Id("p").Dot(srcName)),
					).Else().Block(
						jen.Id(vname).Dot("NULL").Op("=").Op("&").Id("nullBoolTrue"),
					),
				)
				d[fieldName] = jen.Id(vname)
			case avt_string_set, avt_number_set, avt_byte_set:
				arrT := jen.Op("[]*").Id("string")
				if avt == avt_byte_set {
					arrT = jen.Op("[][]").Id("byte")
				}
				stmts = append(stmts,
					jen.Id(arrname).Op(":=").Make(
						arrT,
						jen.Lit(0),
						jen.Len(jen.Id("p").Dot(srcName)),
					),
				)
				needNullBoolTrue = true
				setType := ""
				switch avt {
				case avt_number_set:
					setType = "NS"
					fmtCall := numberFormatStatement(pt, jen.Id(arrix))
					stmts = append(stmts,
						jen.For(jen.List(jen.Id("_"), jen.Id(arrix)).Op(":=").Range().Id("p").Dot(srcName)).Block(
							jen.Id(arrname).Op("=").Append(
								jen.Id(arrname),
								jen.Qual(awsPkg, "String").Call(
									fmtCall,
								),
							),
						),
					)
				case avt_string_set:
					setType = "SS"
					stmts = append(stmts,
						jen.For(jen.List(jen.Id("_"), jen.Id(arrix)).Op(":=").Range().Id("p").Dot(srcName)).Block(
							jen.Id(arrname).Op("=").Append(
								jen.Id(arrname),
								jen.Qual(awsPkg, "String").Call(
									jen.Id(arrix),
								),
							),
						),
					)
				case avt_byte_set:
					setType = "BS"
					stmts = append(stmts,
						jen.For(jen.List(jen.Id("_"), jen.Id(arrix)).Op(":=").Range().Id("p").Dot(srcName)).Block(
							jen.Id(arrname).Op("=").Append(
								jen.Id(arrname),
								jen.Id(arrix),
							),
						),
					)
				}

				stmts = append(stmts,
					jen.Var().Id(vname).Op("*").Qual(dynamoPkg, "AttributeValue"),
					jen.If(jen.Len(jen.Id(arrname)).Op("!=").Lit(0)).Block(
						jen.Id(vname).Op("=").Op("&").Qual(dynamoPkg, "AttributeValue").Values(jen.Dict{
							jen.Id(setType): jen.Id(arrname),
						}),
					).Else().Block(
						jen.Id(vname).Dot("NULL").Op("=").Op("&").Id("nullBoolTrue"),
					),
				)
				d[fieldName] = jen.Id(vname)
			}
		}

		if needNullBoolTrue {
			stmts = append([]jen.Code{
				jen.Id("nullBoolTrue").Op(":=").True(),
			}, stmts...)
		}

		if needProtoBuffer {
			stmts = append([]jen.Code{
				jen.Id(protoBuffer).Op(":=").Qual(protoPkg, "NewBuffer").Call(jen.Nil()),
			}, stmts...)
		}

		if needErr {
			stmts = append([]jen.Code{
				jen.Op("var").Id("err").Id("error"),
			}, stmts...)
		}

		if needStringBuilder {
			stmts = append([]jen.Code{
				jen.Op("var").Id("sb").Qual(stringsPkg, "Builder"),
			}, stmts...)
		}

		stmts = append(stmts, jen.Id("av").Dot("M").Op("=").Map(jen.String()).Op("*").Qual(dynamoPkg, "AttributeValue").Values(d))

		stmts = append(stmts, jen.Return(jen.Nil()))

		f.Func().Params(
			jen.Id("p").Op("*").Id(structName.String()),
		).Id("MarshalDynamoDBAttributeValue").Params(jen.Id("av").Op("*").Qual(dynamoPkg, "AttributeValue")).Id("error").Block(
			stmts...,
		).Line()
	}
	return nil
}

func (m *Module) applyUnmarshal(f *jen.File, in pgs.File) error {
	for _, msg := range in.AllMessages() {
		structName := m.ctx.Name(msg)
		mext := dynamopb.DynamoMessageOptions{}
		ok, err := msg.Extension(dynamopb.E_Msg, &mext)
		if err != nil {
			m.Logf("Parsing dynamo.msg failed: %s", err)
			m.Fail("code generation failed")
		}
		if ok && mext.Disabled {
			m.Logf("dynamo.msg disabled for %s", structName)
			continue
		}

		stmts := []jen.Code{}

		typeName := fmt.Sprintf("%s.%s", msg.Package().ProtoName().String(), msg.Name())

		stmts = append(stmts,
			jen.List(jen.Id(typeField), jen.Id("ok")).Op(":=").Id("av").Dot("M").Index(jen.Lit(typeField)),
			jen.If(jen.Op("!").Id("ok")).Block(
				jen.Return(jen.Qual(fmtPkg, "Errorf").Call(
					jen.Lit("dyanmo: "+typeField+" missing"),
				),
				)),
			jen.If(jen.Qual(awsPkg, "StringValue").Call(jen.Id(typeField).Dot("S")).Op("!=").Lit(typeName)).Block(
				jen.Return(jen.Qual(fmtPkg, "Errorf").Call(
					jen.Lit(fmt.Sprintf("dyanmo: _type mismatch: %s expected, got: '%s'", typeName, "%s")),
					jen.Id(typeField),
				),
				)),
		)

		stmts = append(stmts,
			jen.List(jen.Id(valueField), jen.Id("ok")).Op(":=").Id("av").Dot("M").Index(jen.Lit(valueField)),
			jen.If(jen.Op("!").Id("ok")).Block(
				jen.Return(jen.Qual(fmtPkg, "Errorf").Call(
					jen.Lit("dyanmo: "+valueField+" missing"),
				),
				)),
			jen.Return(jen.Qual(protoPkg, "Unmarshal").Call(jen.Id(valueField).Dot("B"), jen.Id("p"))),
		)

		f.Func().Params(
			jen.Id("p").Op("*").Id(structName.String()),
		).Id("UnmarshalDynamoDBAttributeValue").Params(jen.Id("av").Op("*").Qual(dynamoPkg, "AttributeValue")).Id("error").Block(
			stmts...,
		).Line()

		f.Func().Params(
			jen.Id("p").Op("*").Id(structName.String()),
		).Id("UnmarshalDynamo").Params(jen.Id("av").Op("*").Qual(dynamoPkg, "AttributeValue")).Id("error").Block(
			jen.Return(jen.Id("p").Dot("UnmarshalDynamoDBAttributeValue").Call(jen.Id("av"))),
		).Line()

		f.Func().Params(
			jen.Id("p").Op("*").Id(structName.String()),
		).Id("UnmarshalDynamoItem").Params(jen.Id("av").Map(jen.String()).Op("*").Qual(dynamoPkg, "AttributeValue")).Id("error").Block(
			jen.Return(jen.Id("p").Dot("UnmarshalDynamoDBAttributeValue").Call(
				jen.Op("&").Qual(dynamoPkg, "AttributeValue").Values(jen.Dict{
					jen.Id("M"): jen.Id("av"),
				}),
			)),
		).Line()
	}

	return nil
}

func numberFormatStatement(pt pgs.ProtoType, access *jen.Statement) *jen.Statement {
	var rv *jen.Statement
	switch pt {
	case pgs.DoubleT, pgs.FloatT:
		rv = jen.Qual(strconvPkg, "FormatFloat").Call(
			jen.Id("float64").Call(access),
			jen.LitByte('E'),
			jen.Lit(-1),
			jen.Lit(64),
		)
	case pgs.Int64T, pgs.SFixed64, pgs.SInt64, pgs.Int32T, pgs.SFixed32, pgs.SInt32, pgs.EnumT:
		rv = jen.Qual(strconvPkg, "FormatInt").Call(
			jen.Id("int64").Call(access),
			jen.Lit(10),
		)
	case pgs.UInt64T, pgs.Fixed64T, pgs.UInt32T, pgs.Fixed32T:
		rv = jen.Qual(strconvPkg, "FormatUint").Call(
			jen.Id("uint64").Call(access),
			jen.Lit(10),
		)
	}
	return rv
}

func numberParseStatement(pt pgs.ProtoType, access *jen.Statement) *jen.Statement {
	var rv *jen.Statement
	switch pt {
	case pgs.DoubleT, pgs.FloatT:
		rv = jen.Qual(strconvPkg, "ParseFloat").Call(
			jen.Qual(awsPkg, "StringValue").Call(access),
			jen.Lit(64),
		)
	case pgs.Int64T, pgs.SFixed64, pgs.SInt64, pgs.Int32T, pgs.SFixed32, pgs.SInt32:
		rv = jen.Qual(strconvPkg, "ParseInt").Call(
			jen.Qual(awsPkg, "StringValue").Call(access),
			jen.Lit(10),
			jen.Lit(64),
		)
	case pgs.UInt64T, pgs.Fixed64T, pgs.UInt32T, pgs.Fixed32T:
		rv = jen.Qual(strconvPkg, "ParseUint").Call(
			jen.Qual(awsPkg, "StringValue").Call(access),
			jen.Lit(10),
			jen.Lit(64),
		)
	}
	return rv
}
