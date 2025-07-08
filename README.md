NVM (Node Version Manager) em Go
Um gerenciador de versões do Node.js escrito em Go, inspirado no NVM original.
Funcionalidades
✅ Instalar versões do Node.js
✅ Alternar entre versões instaladas
✅ Listar versões instaladas
✅ Listar versões disponíveis remotamente
✅ Desinstalar versões
✅ Mostrar versão atual
✅ Mostrar caminho do executável
✅ Auto-complete para bash
Instalação
Pré-requisitos
Go 1.16 ou superior
Sistema operacional: Linux, macOS ou Windows
Instalação Rápida
bash
# Clone ou baixe os arquivos
git clone <seu-repositorio>
cd nvm-go

# Compile e instale
make install

# Adicione ao seu shell
echo 'source ~/.nvm/nvm.sh' >> ~/.bashrc
source ~/.bashrc
Instalação Manual
bash
# Compile o binário
go build -o nvm main.go

# Crie o diretório
mkdir -p ~/.nvm

# Copie os arquivos
cp nvm ~/.nvm/
cp nvm.sh ~/.nvm/

# Adicione ao seu shell
echo 'source ~/.nvm/nvm.sh' >> ~/.bashrc
source ~/.bashrc
Uso
Comandos Básicos
bash
# Instalar uma versão
nvm install 18.17.0
nvm install 20.5.0

# Usar uma versão
nvm use 18.17.0

# Listar versões instaladas
nvm list

# Listar versões disponíveis (primeiras 20)
nvm list-remote

# Verificar versão atual
nvm current

# Mostrar caminho do executável
nvm which

# Desinstalar uma versão
nvm uninstall 18.17.0

# Ajuda
nvm help
Exemplos de Uso
bash
# Instalar a versão LTS mais recente
nvm install 20.5.0

# Usar a versão instalada
nvm use 20.5.0

# Verificar se está funcionando
node --version
npm --version

# Instalar outra versão
nvm install 18.17.0

# Alternar entre versões
nvm use 18.17.0
nvm use 20.5.0

# Ver todas as versões instaladas
nvm list
Estrutura de Diretórios
~/.nvm/
├── nvm              # Binário principal
├── nvm.sh           # Script shell
├── current/         # Link simbólico para versão atual
└── versions/        # Versões instaladas
    ├── v18.17.0/
    │   ├── bin/
    │   │   ├── node
    │   │   └── npm
    │   └── ...
    └── v20.5.0/
        ├── bin/
        │   ├── node
        │   └── npm
        └── ...
Diferenças do NVM Original
Similaridades
Comandos principais (install, use, list, etc.)
Estrutura de diretórios similar
Funcionalidade de links simbólicos
Diferenças
Escrito em Go (mais rápido)
Binário único (não depende de scripts bash complexos)
Suporte nativo a múltiplas plataformas
Download e extração mais eficientes
Desenvolvimento
Compilar
bash
make build
Instalar localmente
bash
make install
Limpar
bash
make clean
Desinstalar
bash
make uninstall
Executar testes
bash
make test
Arquitetura
O projeto é composto por:
main.go - Programa principal em Go
nvm.sh - Script shell para integração com o terminal
Makefile - Automação de build e instalação
Componentes Principais
NodeVersionManager: Struct principal que gerencia versões
Download: Sistema de download direto do nodejs.org
Extração: Suporte a arquivos .tar.gz
Links Simbólicos: Gerenciamento de versão atual
Shell Integration: Integração com bash/zsh
Limitações Conhecidas
Não suporte a arquivos .nvmrc (pode ser adicionado)
Não migra configurações do NVM original automaticamente
Auto-complete limitado para versões remotas (performance)
Contribuindo
Faça um fork do projeto
Crie uma branch para sua feature
Faça commit das mudanças
Faça push para a branch
Abra um Pull Request
Licença
MIT License - veja o arquivo LICENSE para detalhes.
Roadmap
 Suporte a arquivos .nvmrc
 Migração automática do NVM original
 Suporte a mais formatos de arquivo
 Cache de versões remotas
 Testes automatizados
 CI/CD
 Instalação via package managers
Troubleshooting
Problema: "nvm: command not found"
Solução: Verifique se adicionou source ~/.nvm/nvm.sh ao seu ~/.bashrc
Problema: "Permission denied"
Solução: Execute chmod +x ~/.nvm/nvm
Problema: "Version not found"
Solução: Verifique se a versão existe com nvm list-remote
Problema: PATH não atualizado
Solução: Feche e abra o terminal, ou execute source ~/.bashrc
