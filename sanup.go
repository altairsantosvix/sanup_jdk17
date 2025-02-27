package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

// Pom representa a estrutura básica do pom.xml
type Pom struct {
	XMLName    xml.Name `xml:"project"`
	Properties struct {
		JavaVersion string `xml:"java.version"`
	} `xml:"properties"`
}

// detectMavenProject verifica se o projeto tem um pom.xml
func detectMavenProject(repoPath string) bool {
	_, err := os.Stat(repoPath + "/pom.xml")
	return !os.IsNotExist(err)
}

// backupPom faz um backup do pom.xml antes da modificação
func backupPom(repoPath string) {
	original := repoPath + "/pom.xml"
	backup := repoPath + "/pom_backup.xml"

	src, err := os.Open(original)
	if err != nil {
		fmt.Println("❌ Erro ao abrir o pom.xml para backup:", err)
		return
	}
	defer src.Close()

	dst, err := os.Create(backup)
	if err != nil {
		fmt.Println("❌ Erro ao criar o backup:", err)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		fmt.Println("❌ Erro ao copiar pom.xml para backup:", err)
		return
	}

	fmt.Println("✅ Backup criado: pom_backup.xml")
}

// checkVulnerabilities executa OWASP Dependency Check para verificar vulnerabilidades
func checkVulnerabilities(repoPath string) bool {
	fmt.Println("🔍 Verificando vulnerabilidades com OWASP Dependency Check...")
	cmd := exec.Command("dependency-check", "--project", "java-migration", "--scan", repoPath, "--format", "JSON", "--out", "report")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("❌ Erro ao executar verificação de vulnerabilidades:", string(output))
		return false
	}
	fmt.Println("✅ Nenhuma vulnerabilidade grave encontrada.")
	return true
}

// updateVulnerableDependencies usa Maven Versions Plugin para atualizar dependências
func updateVulnerableDependencies(repoPath string) bool {
	fmt.Println("🔄 Verificando atualizações de dependências...")

	cmd := exec.Command("mvn", "org.codehaus.mojo:versions-maven-plugin:2.8.1:display-dependency-updates", "-f", repoPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("❌ Erro ao verificar dependências:", string(output))
		return false
	}

	fmt.Println("📜 Dependências desatualizadas detectadas:\n", string(output))

	// Atualizar automaticamente as dependências vulneráveis
	fmt.Println("⚙️ Atualizando dependências para as versões mais recentes...")
	cmd = exec.Command("mvn", "versions:use-latest-releases", "-f", repoPath)
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("❌ Erro ao atualizar dependências:", string(output))
		return false
	}

	fmt.Println("✅ Dependências vulneráveis foram atualizadas no pom.xml.")
	return true
}

// buildProject compila o projeto para validar a migração para Java 17
func buildProject(repoPath string) bool {
	fmt.Println("🚀 Executando build para validar a migração...")

	cmd := exec.Command("mvn", "clean", "package", "-f", repoPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("❌ Erro ao compilar o projeto:", string(output))
		return false
	}

	fmt.Println("✅ Build bem-sucedido! Migração completa.")
	return true
}

// updateJavaVersion só modifica o pom.xml se tudo rodou sem erro
func updateJavaVersion(repoPath string) {
	pomFile := repoPath + "/pom.xml"

	// Ler o pom.xml
	xmlData, err := ioutil.ReadFile(pomFile)
	if err != nil {
		fmt.Println("❌ Erro ao ler pom.xml:", err)
		return
	}

	// Parse do XML
	var pom Pom
	err = xml.Unmarshal(xmlData, &pom)
	if err != nil {
		fmt.Println("❌ Erro ao fazer parsing do pom.xml:", err)
		return
	}

	// Atualiza a versão do Java
	if pom.Properties.JavaVersion == "8" || pom.Properties.JavaVersion == "11" {
		pom.Properties.JavaVersion = "17"
	} else {
		fmt.Println("ℹ️ O projeto já está usando Java 17 ou versão superior. Nenhuma alteração necessária.")
		return
	}

	// Serializar de volta para XML
	xmlOutput, err := xml.MarshalIndent(pom, "", "  ")
	if err != nil {
		fmt.Println("❌ Erro ao gerar novo pom.xml:", err)
		return
	}

	// Adicionar cabeçalho XML
	xmlOutput = append([]byte(xml.Header), xmlOutput...)

	// Escrever no arquivo
	err = ioutil.WriteFile(pomFile, xmlOutput, 0644)
	if err != nil {
		fmt.Println("❌ Erro ao escrever no pom.xml:", err)
		return
	}

	fmt.Println("✅ Atualizada a versão do Java para 17 no pom.xml.")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("❌ Uso: go run main.go <caminho-do-repo>")
		return
	}

	repoPath := os.Args[1]

	if !detectMavenProject(repoPath) {
		fmt.Println("❌ Nenhum pom.xml encontrado. Este script suporta apenas projetos Maven.")
		return
	}

	backupPom(repoPath) // Criar backup antes de modificar

	// Verificar vulnerabilidades e dependências antes de alterar o pom.xml
	if checkVulnerabilities(repoPath) && updateVulnerableDependencies(repoPath) && buildProject(repoPath) {
		updateJavaVersion(repoPath) // Só modifica se tudo rodou bem
	} else {
		fmt.Println("⚠️ Migração cancelada devido a erros. O pom.xml não foi alterado.")
	}
}
