package compiler

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/WebDelve/activeledger-contract-compiler/config"
	"github.com/WebDelve/activeledger-contract-compiler/helper"
)

type Compiler struct {
	config        *config.Config
	contractEntry string
	contractDir   string
	contracts     map[string]Contract
	fileList      []string
	toProcess     []string
	compiled      []string
}

type Contract struct {
	name            string
	imports         []string
	externalImports []string
	localImports    []string
	body            []string
}

type importData struct {
	packageName string
	classes     []string
}

func GetCompiler(config *config.Config, contractEntry string) Compiler {
	contract := make(map[string]Contract)

	lastSlash := strings.LastIndex(contractEntry, "/")

	contractDir := contractEntry[:lastSlash]

	return Compiler{
		config,
		contractEntry,
		contractDir,
		contract,
		[]string{},
		[]string{},
		[]string{},
	}
}

func (c *Compiler) Compile() {
	entryLines := readLines(c.contractEntry)
	c.fileList = getFileNamesInDir(c.contractEntry)

	fileIndex := strings.LastIndex(c.contractEntry, "/")
	entryName := cleanFileName(c.contractEntry[fileIndex+1:])
	c.buildContract(entryName, entryLines)

	for _, name := range c.fileList {
		if name != c.contractEntry {
			c.toProcess = append(c.toProcess, name)
		}
	}

	for _, f := range c.toProcess {
		name := f
		path := fmt.Sprintf("%s/%s.ts", c.contractDir, name)
		lines := readLines(path)

		c.buildContract(name, lines)
	}

	c.combine()

	writeToFile(c.config.Output, c.compiled)
}

func (c *Compiler) combine() {

	imports := []string{}
	body := []string{}

	for _, contract := range c.contracts {
		imports = append(imports, contract.externalImports...)

		body = append(body, fmt.Sprintf("// Included from file %s.ts", contract.name))
		body = append(body, contract.body...)
		body = append(body, "")
	}

	imports = buildImports(imports)

	output := []string{}

	output = append(output, imports...)
	output = append(output, "")
	output = append(output, body...)

	c.compiled = output

}

func (c *Compiler) buildContract(name string, lines []string) {
	contract := Contract{
		name:            name,
		imports:         []string{},
		externalImports: []string{},
		localImports:    []string{},
		body:            []string{},
	}

	for _, line := range lines {
		reg, err := regexp.Compile(`import+\s*{\s*([A-z]+,*\s*)*}\s*from\s*(\"|\').*`)
		if err != nil {
			helper.Error(err, "Error checking line against regex")
		}

		if reg.MatchString(line) {
			importName, isLocal := c.processImport(line)
			if isLocal {
				contract.localImports = append(contract.localImports, importName)
			} else {
				contract.externalImports = append(contract.externalImports, line)
			}

			contract.imports = append(contract.imports, line)
			continue
		}

		contract.body = append(contract.body, line)
	}

	c.contracts[name] = contract
}

func (c *Compiler) processImport(line string) (string, bool) {
	quoteIndexes := []int{}

	for i := len(line) - 1; i >= 0; i-- {
		if line[i] == '"' {
			quoteIndexes = append(quoteIndexes, i)
		}

		if len(quoteIndexes) >= 2 {
			break
		}
	}

	importFile := line[quoteIndexes[1]+1 : quoteIndexes[0]]

	importSplit := strings.Split(importFile, "/")
	clean := importSplit[len(importSplit)-1]

	isLocal := false
	for _, v := range c.fileList {
		if v == clean {
			isLocal = true
			break
		}
	}

	return clean, isLocal
}

func getFileNamesInDir(path string) []string {
	fileList := []string{}

	fileIndex := strings.LastIndex(path, "/")
	pathToRead := path[0:fileIndex]

	files, err := os.ReadDir(pathToRead)
	if err != nil {
		helper.Error(err, "Error reading directory")
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := file.Name()

		// If the name ends with .ts
		if strings.Contains(name[len(name)-3:], ".ts") {
			clean := cleanFileName(name)
			fileList = append(fileList, clean)
		}

	}

	return fileList
}

func cleanFileName(raw string) string {
	extIndex := strings.LastIndex(raw, ".ts")
	clean := raw[0:extIndex]

	return clean
}

func getClasses(line string) []string {
	str := strings.Index(line, "{") + 1
	end := strings.Index(line, "}")

	classString := line[str:end]

	if strings.Contains(classString, ",") {
		return strings.Split(classString, ",")
	}

	return []string{classString}
}

func getPackage(line string) string {
	str := strings.Index(line, "\"") + 1
	end := strings.LastIndex(line, "\"")

	return line[str:end]
}

func buildImports(imports []string) []string {
	cleaned := []string{}
	importCache := make(map[string]importData)

	for _, i := range imports {
		iData := importData{}

		iData.classes = getClasses(i)

		packageName := getPackage(i)
		iData.packageName = packageName

		// check if import cache already has importdata stored under this key
		if len(importCache[packageName].classes) > 0 {

			if entry, ok := importCache[packageName]; ok {
				storedClasses := entry.classes
				// merge non duplicates
				for _, cl := range iData.classes {
					if contains(storedClasses, cl) {
						continue
					}

					entry.classes = append(entry.classes, cl)
				}
			}
		}

		importCache[packageName] = iData

	}

	for k, v := range importCache {
		impStr := buildImportLine(k, v.classes)
		cleaned = append(cleaned, impStr)
	}

	return cleaned
}

func buildImportLine(packageName string, classes []string) string {
	importLine := "import {"

	for i := 0; i < len(classes); i++ {
		importLine = importLine + classes[i]

		if i != len(classes)-1 {
			importLine = importLine + ","
		}
	}

	importLine = fmt.Sprintf("%s} from \"%s\";", importLine, packageName)

	return importLine
}

func contains(sl []string, str string) bool {
	for _, v := range sl {
		if v == str {
			return true
		}
	}

	return false
}

func readLines(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		helper.Error(err, fmt.Sprintf("Error opening contract files at %s", path))
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := []string{}

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		helper.Error(err, fmt.Sprintf("Error building line array for contract %s", path))
	}

	return lines
}

func writeToFile(path string, lines []string) {
	f, err := os.OpenFile(path, os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0644)

	if err != nil {
		helper.Error(err, "Error opening file to write")
	}

	defer f.Close()

	writer := bufio.NewWriter(f)

	for _, line := range lines {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			helper.Error(err, "Error writing line to compiled file")
		}
	}

	writer.Flush()
}
