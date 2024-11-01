package generator_test

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/Djarvur/protoc-gen-python-grpc/internal/generator"
	"github.com/Djarvur/protoc-gen-python-grpc/internal/template"
)

const (
	fileNameSuffixDefault = "_pb2_grpc.py"
)

type GeneratorSuite struct {
	suite.Suite
	InBytes  []byte
	OutBytes []byte
}

func TestGeneratorSuite(t *testing.T) {
	t.Parallel()

	// read serialized request/response from testdata
	in, err := os.OpenFile("testdata/in.bin", os.O_RDONLY, 0)
	if err != nil {
		t.Fatalf("not found testdata/in.bin")
	}
	defer in.Close()

	inBytes, err := io.ReadAll(in)
	if err != nil {
		t.Fatalf("error reading testdata/in.bin")
	}

	out, err := os.OpenFile("testdata/out.bin", os.O_RDONLY, 0)
	if err != nil {
		t.Fatalf("not found testdata/out.bin")
	}
	defer out.Close()

	outBytes, err := io.ReadAll(out)
	if err != nil {
		t.Fatalf("error reading testdata/out.bin")
	}

	// create suite and set test data
	s := new(GeneratorSuite)
	s.InBytes = inBytes
	s.OutBytes = outBytes

	suite.Run(t, s)
}

func (s *GeneratorSuite) TestGenerator_New() {
	// Arrange
	suffix := "suffix"
	tmplSrc := "template"

	// Act
	gen := generator.New(suffix, tmplSrc)

	// Assert
	s.Require().NotNil(gen, "got nil generator")
	s.Equal(suffix, gen.Suffix)
	s.Equal(tmplSrc, gen.Template)
}

func (s *GeneratorSuite) TestGenerator_Generate_Success() {
	// Arrange
	req := new(pluginpb.CodeGeneratorRequest)
	if err := proto.Unmarshal(s.InBytes, req); err != nil {
		s.T().Fatalf("error unmarshalling testdata/in.bin: %v", err)
	}

	templateSource := template.NewTemplateValue()
	gen := generator.New(fileNameSuffixDefault, templateSource.Source())

	// Act
	resp, err := gen.Generate(req)

	// Assert
	s.Require().NoError(err, "error generating response")
	s.Require().NotNil(resp, "got nil response")

	outData, err := proto.Marshal(resp)
	if err != nil {
		s.T().Fatalf("error marshalling generated response: %v", err)
	}

	s.Require().Equal(s.OutBytes, outData)
}

func (s *GeneratorSuite) TestGenerator_Generate_EmptyRequestSuccess() {
	// Arrange
	req := new(pluginpb.CodeGeneratorRequest)

	expectedResp := new(pluginpb.CodeGeneratorResponse)
	expectedResp.SupportedFeatures = proto.Uint64(generator.SupportedFeatures)

	expectedData, err := proto.Marshal(expectedResp)
	if err != nil {
		s.T().Fatalf("error marshaling expected response: %v", err)
	}

	gen := generator.New("", "")

	// Act
	resp, err := gen.Generate(req)

	// Assert
	s.Require().NoError(err, "error generating response")
	s.Require().NotNil(resp, "got nil response")

	outData, err := proto.Marshal(resp)
	if err != nil {
		s.T().Fatalf("error marshalling generated response: %v", err)
	}

	s.Require().Equal(expectedData, outData)
}

func (s *GeneratorSuite) TestGenerator_Generate_BuildTemplateError() {
	// Arrange
	req := new(pluginpb.CodeGeneratorRequest)
	gen := generator.New("", "{{ if }}")

	// Act
	_, err := gen.Generate(req)

	// Assert
	s.Require().Error(err)
	s.Require().ErrorIs(err, generator.ErrTemplateBuild)
}

func (s *GeneratorSuite) TestGenerator_Generate_ExecuteTemplateError() {
	// Arrange
	req := new(pluginpb.CodeGeneratorRequest)
	if err := proto.Unmarshal(s.InBytes, req); err != nil {
		s.T().Errorf("error unmarshalling testdata/in.bin: %v", err)
	}

	gen := generator.New(fileNameSuffixDefault, "{{ .Data }}")

	// Act
	_, err := gen.Generate(req)

	// Assert
	s.Require().Error(err)
	s.Require().ErrorIs(err, generator.ErrTemplateExec)
}
