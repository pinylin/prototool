// Copyright (c) 2019 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package lint

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/emicklei/proto"
	"github.com/uber/prototool/internal/file"
	"github.com/uber/prototool/internal/settings"
	"github.com/uber/prototool/internal/text"
	"go.uber.org/zap"
)

var (
	// AllLinters is the slice of all known Linters.
	AllLinters = []Linter{
		commentsNoCStyleLinter,
		commentsNoInlineLinter,
		enumFieldNamesUppercaseLinter,
		enumFieldNamesUpperSnakeCaseLinter,
		enumFieldPrefixesLinter,
		enumFieldPrefixesExceptMessageLinter,
		enumNamesCamelCaseLinter,
		enumNamesCapitalizedLinter,
		enumZeroValuesInvalidLinter,
		enumZeroValuesInvalidExceptMessageLinter,
		enumsHaveCommentsLinter,
		enumsHaveSentenceCommentsLinter,
		enumsNoAllowAliasLinter,
		fieldsNotReservedLinter,
		fileHeaderLinter,
		fileOptionsCSharpNamespaceSameInDirLinter,
		fileOptionsEqualCSharpNamespaceCapitalizedLinter,
		fileOptionsEqualGoPackageV2SuffixLinter,
		fileOptionsEqualGoPackagePbSuffixLinter,
		fileOptionsEqualJavaMultipleFilesTrueLinter,
		fileOptionsEqualJavaOuterClassnameProtoSuffixLinter,
		fileOptionsEqualJavaPackageComPrefixLinter,
		fileOptionsEqualOBJCClassPrefixAbbrLinter,
		fileOptionsEqualPHPNamespaceCapitalizedLinter,
		fileOptionsGoPackageNotLongFormLinter,
		fileOptionsGoPackageSameInDirLinter,
		fileOptionsJavaMultipleFilesSameInDirLinter,
		fileOptionsJavaPackageSameInDirLinter,
		fileOptionsOBJCClassPrefixSameInDirLinter,
		fileOptionsPHPNamespaceSameInDirLinter,
		fileOptionsRequireCSharpNamespaceLinter,
		fileOptionsRequireGoPackageLinter,
		fileOptionsRequireJavaMultipleFilesLinter,
		fileOptionsRequireJavaOuterClassnameLinter,
		fileOptionsRequireJavaPackageLinter,
		fileOptionsRequireOBJCClassPrefixLinter,
		fileOptionsRequirePHPNamespaceLinter,
		fileOptionsUnsetJavaMultipleFilesLinter,
		fileOptionsUnsetJavaOuterClassnameLinter,
		gogoNotImportedLinter,
		importsNotPublicLinter,
		importsNotWeakLinter,
		messageFieldsNotFloatsLinter,
		messageFieldsNoJSONNameLinter,
		messageFieldNamesLowerSnakeCaseLinter,
		messageFieldNamesLowercaseLinter,
		messageNamesCamelCaseLinter,
		messageNamesCapitalizedLinter,
		messagesHaveCommentsLinter,
		messagesHaveCommentsExceptRequestResponseTypesLinter,
		messagesHaveSentenceCommentsExceptRequestResponseTypesLinter,
		namesNoCommonLinter,
		namesNoDataLinter,
		namesNoUUIDLinter,
		oneofNamesLowerSnakeCaseLinter,
		packageIsDeclaredLinter,
		packageLowerSnakeCaseLinter,
		packageMajorBetaVersionedLinter,
		packageNoKeywordsLinter,
		packagesSameInDirLinter,
		rpcsHaveCommentsLinter,
		rpcsHaveSentenceCommentsLinter,
		rpcNamesCamelCaseLinter,
		rpcNamesCapitalizedLinter,
		rpcOptionsNoGoogleAPIHTTPLinter,
		requestResponseTypesAfterServiceLinter,
		requestResponseTypesInSameFileLinter,
		requestResponseTypesOnlyInFileLinter,
		requestResponseTypesUniqueLinter,
		requestResponseNamesMatchRPCLinter,
		servicesHaveCommentsLinter,
		servicesHaveSentenceCommentsLinter,
		serviceNamesAPISuffixLinter,
		serviceNamesCamelCaseLinter,
		serviceNamesCapitalizedLinter,
		serviceNamesMatchFileNameLinter,
		serviceNamesNoPluralsLinter,
		syntaxProto3Linter,
		wktDirectlyImportedLinter,
		wktDurationSuffixLinter,
		wktTimestampSuffixLinter,
	}

	// DefaultLinters is the slice of default Linters.
	DefaultLinters = UberLinters

	// GoogleLinters is the slice of linters for the google lint group.
	GoogleLinters = []Linter{
		enumFieldNamesUpperSnakeCaseLinter,
		enumNamesCamelCaseLinter,
		enumNamesCapitalizedLinter,
		fileHeaderLinter,
		messageFieldNamesLowerSnakeCaseLinter,
		messageNamesCamelCaseLinter,
		messageNamesCapitalizedLinter,
		rpcNamesCamelCaseLinter,
		rpcNamesCapitalizedLinter,
		serviceNamesCamelCaseLinter,
		serviceNamesCapitalizedLinter,
	}

	// UberLinters is the slice of linters for the uber lint group.
	UberLinters = []Linter{
		commentsNoCStyleLinter,
		enumFieldNamesUpperSnakeCaseLinter,
		enumFieldPrefixesLinter,
		enumNamesCamelCaseLinter,
		enumNamesCapitalizedLinter,
		enumZeroValuesInvalidLinter,
		enumsNoAllowAliasLinter,
		fileHeaderLinter,
		fileOptionsEqualGoPackagePbSuffixLinter,
		fileOptionsEqualJavaMultipleFilesTrueLinter,
		fileOptionsEqualJavaOuterClassnameProtoSuffixLinter,
		fileOptionsEqualJavaPackageComPrefixLinter,
		fileOptionsGoPackageNotLongFormLinter,
		fileOptionsGoPackageSameInDirLinter,
		fileOptionsJavaMultipleFilesSameInDirLinter,
		fileOptionsJavaPackageSameInDirLinter,
		fileOptionsRequireGoPackageLinter,
		fileOptionsRequireJavaMultipleFilesLinter,
		fileOptionsRequireJavaOuterClassnameLinter,
		fileOptionsRequireJavaPackageLinter,
		messageFieldNamesLowerSnakeCaseLinter,
		messageNamesCamelCaseLinter,
		messageNamesCapitalizedLinter,
		oneofNamesLowerSnakeCaseLinter,
		packageIsDeclaredLinter,
		packageLowerSnakeCaseLinter,
		packagesSameInDirLinter,
		rpcNamesCamelCaseLinter,
		rpcNamesCapitalizedLinter,
		requestResponseTypesInSameFileLinter,
		requestResponseTypesUniqueLinter,
		serviceNamesCamelCaseLinter,
		serviceNamesCapitalizedLinter,
		syntaxProto3Linter,
		wktDirectlyImportedLinter,
	}

	// Uber2Linters is the slice of linters for the uber2 lint group.
	Uber2Linters = []Linter{
		commentsNoCStyleLinter,
		commentsNoInlineLinter,
		enumFieldNamesUpperSnakeCaseLinter,
		enumFieldPrefixesExceptMessageLinter,
		enumNamesCamelCaseLinter,
		enumNamesCapitalizedLinter,
		enumZeroValuesInvalidExceptMessageLinter,
		enumsHaveSentenceCommentsLinter,
		enumsNoAllowAliasLinter,
		fieldsNotReservedLinter,
		fileHeaderLinter,
		fileOptionsCSharpNamespaceSameInDirLinter,
		fileOptionsEqualCSharpNamespaceCapitalizedLinter,
		fileOptionsEqualGoPackageV2SuffixLinter,
		fileOptionsEqualJavaMultipleFilesTrueLinter,
		fileOptionsEqualJavaOuterClassnameProtoSuffixLinter,
		fileOptionsEqualJavaPackageComPrefixLinter,
		fileOptionsEqualOBJCClassPrefixAbbrLinter,
		fileOptionsEqualPHPNamespaceCapitalizedLinter,
		fileOptionsGoPackageNotLongFormLinter,
		fileOptionsGoPackageSameInDirLinter,
		fileOptionsJavaMultipleFilesSameInDirLinter,
		fileOptionsJavaPackageSameInDirLinter,
		fileOptionsOBJCClassPrefixSameInDirLinter,
		fileOptionsPHPNamespaceSameInDirLinter,
		fileOptionsRequireCSharpNamespaceLinter,
		fileOptionsRequireGoPackageLinter,
		fileOptionsRequireJavaMultipleFilesLinter,
		fileOptionsRequireJavaOuterClassnameLinter,
		fileOptionsRequireJavaPackageLinter,
		fileOptionsRequireOBJCClassPrefixLinter,
		fileOptionsRequirePHPNamespaceLinter,
		gogoNotImportedLinter,
		importsNotPublicLinter,
		importsNotWeakLinter,
		messagesHaveSentenceCommentsExceptRequestResponseTypesLinter,
		messageFieldNamesLowerSnakeCaseLinter,
		messageFieldsNoJSONNameLinter,
		messageFieldsNotFloatsLinter,
		messageNamesCamelCaseLinter,
		messageNamesCapitalizedLinter,
		namesNoCommonLinter,
		namesNoDataLinter,
		namesNoUUIDLinter,
		oneofNamesLowerSnakeCaseLinter,
		packageIsDeclaredLinter,
		packageLowerSnakeCaseLinter,
		packageMajorBetaVersionedLinter,
		packageNoKeywordsLinter,
		packagesSameInDirLinter,
		rpcsHaveSentenceCommentsLinter,
		rpcNamesCamelCaseLinter,
		rpcNamesCapitalizedLinter,
		rpcOptionsNoGoogleAPIHTTPLinter,
		requestResponseNamesMatchRPCLinter,
		requestResponseTypesAfterServiceLinter,
		requestResponseTypesInSameFileLinter,
		requestResponseTypesOnlyInFileLinter,
		requestResponseTypesUniqueLinter,
		servicesHaveSentenceCommentsLinter,
		serviceNamesAPISuffixLinter,
		serviceNamesCamelCaseLinter,
		serviceNamesCapitalizedLinter,
		serviceNamesMatchFileNameLinter,
		serviceNamesNoPluralsLinter,
		syntaxProto3Linter,
		wktDirectlyImportedLinter,
		wktDurationSuffixLinter,
		wktTimestampSuffixLinter,
	}

	// GroupToLinters is the map from linter group to the corresponding slice of linters.
	GroupToLinters = map[string][]Linter{
		"google": GoogleLinters,
		"uber":   UberLinters,
		"uber2":  Uber2Linters,
	}

	allLintIDs = make(map[string]struct{})
)

func init() {
	for _, linter := range AllLinters {
		if _, ok := allLintIDs[linter.ID()]; ok {
			panic(fmt.Sprintf("duplicate linter id %s", linter.ID()))
		}
		allLintIDs[linter.ID()] = struct{}{}
	}
}

// Runner runs a lint job.
type Runner interface {
	Run(*file.ProtoSet) ([]*text.Failure, error)
}

// RunnerOption is an option for a new Runner.
type RunnerOption func(*runner)

// RunnerWithLogger returns a RunnerOption that uses the given logger.
//
// The default is to use zap.NewNop().
func RunnerWithLogger(logger *zap.Logger) RunnerOption {
	return func(runner *runner) {
		runner.logger = logger
	}
}

// NewRunner returns a new Runner.
func NewRunner(options ...RunnerOption) Runner {
	return newRunner(options...)
}

// FileDescriptor is a wrapper for proto.Proto.
type FileDescriptor struct {
	*proto.Proto

	ProtoSet *file.ProtoSet
	FileData string
}

// The below should not be needed in the CLI
// TODO make private

// Linter is a linter for Protobuf files.
type Linter interface {
	// Return the ID of this Linter. This should be all UPPER_SNAKE_CASE.
	ID() string
	// Return the purpose of this Linter. This should be a human-readable string.
	Purpose() string
	// Check the file data for the descriptors in a common directgory.
	// If there is a lint failure, this returns it in the
	// slice and does not return an error. An error is returned if something
	// unexpected happens. Callers should verify the files are compilable
	// before running this.
	Check(dirPath string, descriptors []*FileDescriptor) ([]*text.Failure, error)
}

// NewLinter is a convenience function that returns a new Linter for the
// given parameters, using a function to record failures.
//
// The ID will be upper-cased.
//
// Failures returned from check do not need to set the ID, this will be overwritten.
func NewLinter(id string, purpose string, addCheck func(func(*text.Failure), string, []*FileDescriptor) error) Linter {
	return newBaseLinter(id, purpose, addCheck)
}

// GetLinters returns the Linters for the LintConfig.
//
// The group, if set, is expected to be lower-case.
//
// The configuration is expected to be valid, deduplicated, and all upper-case.
// IncludeIDs and ExcludeIDs MUST NOT have an intersection.
//
// If the config came from the settings package, this is already validated.
func GetLinters(config settings.LintConfig) ([]Linter, error) {
	var linters []Linter
	var ok bool
	if config.Group != "" {
		linters, ok = GroupToLinters[config.Group]
		if !ok {
			return nil, fmt.Errorf("unknown lint group: %s", config.Group)
		}
	} else if !config.NoDefault {
		// we ignore NoDefault if Group is set
		linters = DefaultLinters
	}
	if len(config.IncludeIDs) == 0 && len(config.ExcludeIDs) == 0 {
		return linters, nil
	}

	// Apply the configured linters to the default group.
	linterMap := make(map[string]Linter, len(linters)+len(config.IncludeIDs)-len(config.ExcludeIDs))
	for _, l := range linters {
		linterMap[l.ID()] = l
	}
	if len(config.IncludeIDs) > 0 {
		for _, l := range AllLinters {
			for _, id := range config.IncludeIDs {
				if err := checkLintID(id); err != nil {
					return nil, err
				}
				if l.ID() == id {
					linterMap[id] = l
				}
			}
		}
	}
	for _, excludeID := range config.ExcludeIDs {
		if err := checkLintID(excludeID); err != nil {
			return nil, err
		}
		delete(linterMap, excludeID)
	}
	for ignoreID := range config.IgnoreIDToFilePaths {
		if err := checkLintID(ignoreID); err != nil {
			return nil, err
		}
	}

	result := make([]Linter, 0, len(linterMap))
	for _, l := range linterMap {
		result = append(result, l)
	}
	return result, nil
}

// GetDirPathToDescriptors is a convenience function that gets the
// descriptors for the given ProtoSet.
func GetDirPathToDescriptors(protoSet *file.ProtoSet) (map[string][]*FileDescriptor, error) {
	dirPathToDescriptors := make(map[string][]*FileDescriptor, len(protoSet.DirPathToFiles))
	for dirPath, protoFiles := range protoSet.DirPathToFiles {
		descriptors := make([]*FileDescriptor, len(protoFiles))
		for i, protoFile := range protoFiles {
			file, err := os.Open(protoFile.Path)
			if err != nil {
				return nil, err
			}
			parser := proto.NewParser(file)
			parser.Filename(protoFile.DisplayPath)
			descriptor, err := parser.Parse()
			if err != nil {
				_ = file.Close()
				return nil, err
			}
			if _, err := file.Seek(0, 0); err != nil {
				_ = file.Close()
				return nil, err
			}
			fileData, err := ioutil.ReadAll(file)
			if err != nil {
				_ = file.Close()
				return nil, err
			}
			_ = file.Close()
			descriptors[i] = &FileDescriptor{
				Proto:    descriptor,
				ProtoSet: protoSet,
				FileData: string(fileData),
			}
		}
		dirPathToDescriptors[dirPath] = descriptors
	}
	return dirPathToDescriptors, nil
}

// CheckMultiple is a convenience function that checks multiple linters and multiple descriptors.
func CheckMultiple(linters []Linter, dirPathToDescriptors map[string][]*FileDescriptor, ignoreIDToFilePaths map[string][]string) ([]*text.Failure, error) {
	var allFailures []*text.Failure
	for dirPath, descriptors := range dirPathToDescriptors {
		for _, linter := range linters {
			failures, err := checkOne(linter, dirPath, descriptors, ignoreIDToFilePaths)
			if err != nil {
				return nil, err
			}
			allFailures = append(allFailures, failures...)
		}
	}
	text.SortFailures(allFailures)
	return allFailures, nil
}

func checkOne(linter Linter, dirPath string, descriptors []*FileDescriptor, ignoreIDToFilePaths map[string][]string) ([]*text.Failure, error) {
	filteredDescriptors, err := filterIgnores(linter, descriptors, ignoreIDToFilePaths)
	if err != nil {
		return nil, err
	}
	return linter.Check(dirPath, filteredDescriptors)
}

func filterIgnores(linter Linter, descriptors []*FileDescriptor, ignoreIDToFilePaths map[string][]string) ([]*FileDescriptor, error) {
	var filteredDescriptors []*FileDescriptor
	for _, descriptor := range descriptors {
		ignore, err := shouldIgnore(linter, descriptor, ignoreIDToFilePaths)
		if err != nil {
			return nil, err
		}
		if !ignore {
			filteredDescriptors = append(filteredDescriptors, descriptor)
		}
	}
	return filteredDescriptors, nil
}

func shouldIgnore(linter Linter, descriptor *FileDescriptor, ignoreIDToFilePaths map[string][]string) (bool, error) {
	filePath := descriptor.Filename
	var err error
	if !filepath.IsAbs(filePath) {
		filePath, err = filepath.Abs(filePath)
		if err != nil {
			return false, err
		}
	}
	ignoreFilePaths, ok := ignoreIDToFilePaths[linter.ID()]
	if !ok {
		return false, nil
	}
	for _, ignoreFilePath := range ignoreFilePaths {
		if filePath == ignoreFilePath {
			return true, nil
		}
	}
	return false, nil
}

func checkLintID(lintID string) error {
	if _, ok := allLintIDs[lintID]; !ok {
		return fmt.Errorf("unknown lint id in configuration file: %s", lintID)
	}
	return nil
}

func hasGolangStyleComment(comment *proto.Comment, name string) bool {
	return comment != nil && len(comment.Lines) > 0 && strings.HasPrefix(comment.Lines[0], fmt.Sprintf(" %s ", name))
}

func hasCompleteSentenceComment(comment *proto.Comment) bool {
	return commentStartsWithUppercaseLetter(comment) && commentContainsPeriod(comment)
}

func commentStartsWithUppercaseLetter(comment *proto.Comment) bool {
	if comment == nil || len(comment.Lines) == 0 {
		return false
	}
	firstLine := strings.TrimSpace(comment.Lines[0])
	if firstLine == "" {
		return false
	}
	return unicode.IsUpper(rune(firstLine[0])) || unicode.IsDigit(rune(firstLine[0]))
}

func commentContainsPeriod(comment *proto.Comment) bool {
	if comment == nil || len(comment.Lines) == 0 {
		return false
	}
	// very primitive check, could make better with NLP but this is hard with comments
	// since comments can contain code examples, links, etc.
	for _, line := range comment.Lines {
		if strings.Contains(line, ".") {
			return true
		}
	}
	return false
}

func isSuppressed(comment *proto.Comment, annotation string) bool {
	if comment == nil {
		return false
	}
	annotation = "@suppresswarnings " + annotation
	for _, line := range comment.Lines {
		if strings.Contains(line, annotation) {
			return true
		}
	}
	return false
}
