package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// Pom representa a estrutura básica do pom.xml
type Pom struct {
	XMLName    xml.Name `xml:"project"`
	Properties struct {
		JavaVersion string `xml:"java.version"`
	} `xml:"properties"`
}

func detectMavenProject(repoPath string) bool {
	_, err := os.Stat(repoPath + "/pom.xml")
	return !os.IsNotExist(err)
}

func updateJavaVersion(repoPath string) {
	pomFile := repoPath + "/pom.xml"

	// Ler o pom.xml
	xmlData, err := ioutil.ReadFile(pomFile)
	if err != nil {
		fmt.Println("Erro ao ler pom.xml:", err)
		return
	}

	// Parse do XML
	var pom Pom
	err = xml.Unmarshal(xmlData, &pom)
	if err != nil {
		fmt.Println("Erro ao fazer parsing do pom.xml:", err)
		return
	}

	// Atualiza a versão do Java
	if pom.Properties.JavaVersion == "8" || pom.Properties.JavaVersion == "11" {
		pom.Properties.JavaVersion = "17"
	}

	// Serializar de volta para XML
	xmlOutput, err := xml.MarshalIndent(pom, "", "  ")
	if err != nil {
		fmt.Println("Erro ao gerar novo pom.xml:", err)
		return
	}

	// Adicionar cabeçalho XML
	xmlOutput = append([]byte(xml.Header), xmlOutput...)

	// Escrever no arquivo
	err = ioutil.WriteFile(pomFile, xmlOutput, 0644)
	if err != nil {
		fmt.Println("Erro ao escrever no pom.xml:", err)
		return
	}

	fmt.Println("✅ Atualizada a versão do Java para 17 no pom.xml.")
}

func checkVulnerabilities(repoPath string) {
	fmt.Println("🔍 Verificando vulnerabilidades com OWASP Dependency Check...")
	cmd := exec.Command("dependency-check", "--project", "java-migration", "--scan", repoPath, "--format", "JSON", "--out", "report")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("❌ Erro ao executar verificação de vulnerabilidades:", string(output))
		return
	}
	fmt.Println("✅ Verificação de vulnerabilidades concluída. Veja o relatório em './report/dependency-check-report.json'.")
}

func updateVulnerableDependencies(repoPath string) {
	fmt.Println("🔄 Verificando atualizações de dependências...")

	cmd := exec.Command("mvn", "org.codehaus.mojo:versions-maven-plugin:2.8.1:display-dependency-updates", "-f", repoPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("❌ Erro ao verificar dependências:", string(output))
		return
	}

	fmt.Println("📜 Dependências desatualizadas detectadas:\n", string(output))

	// Atualizar automaticamente as dependências vulneráveis
	fmt.Println("⚙️ Atualizando dependências para as versões mais recentes...")
	cmd = exec.Command("mvn", "versions:use-latest-releases", "-f", repoPath)
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("❌ Erro ao atualizar dependências:", string(output))
		return
	}

	fmt.Println("✅ Dependências vulneráveis foram atualizadas no pom.xml.")
}

func buildProject(repoPath string) {
	fmt.Println("🚀 Executando build para validar a migração...")

	cmd := exec.Command("mvn", "clean", "package", "-f", repoPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("❌ Erro ao compilar o projeto:", string(output))
		return
	}

	fmt.Println("✅ Build bem-sucedido! Migração completa.")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("❌ Uso: go run migrate.go <caminho-do-repo>")
		return
	}

	repoPath := os.Args[1]

	if !detectMavenProject(repoPath) {
		fmt.Println("❌ Nenhum pom.xml encontrado. Este script suporta apenas projetos Maven.")
		return
	}

	updateJavaVersion(repoPath)
	checkVulnerabilities(repoPath)
	updateVulnerableDependencies(repoPath)
	buildProject(repoPath)
}
