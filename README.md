# url_shortener

> Nota: a secção abaixo é uma chuleta de referência de Go (comandos, sintaxe, ecossistema), pensada para consultares e depois apagares. Contrastes com Rust incluídos onde ajudam.

## Comandos essenciais

| Comando | Para quê | Equivalente Rust |
|---|---|---|
| `go mod init <nome>` | Cria `go.mod`, inicia um módulo novo | `cargo init` |
| `go get <path>` | Adiciona/atualiza uma dependência | `cargo add` |
| `go mod tidy` | Remove deps não usadas e adiciona as que faltam nos imports | `cargo add`/limpeza manual, não há 1:1 exato |
| `go build ./...` | Compila tudo, não gera nada se só quiseres validar (`go build -o /dev/null ./...`) | `cargo build` |
| `go run ./cmd/server` | Compila e corre já (binário temporário) | `cargo run` |
| `go install <path>` | Compila e instala o binário em `$GOBIN`/`$GOPATH/bin` | `cargo install` |
| `go test ./...` | Corre todos os testes do módulo | `cargo test` |
| `go test ./... -run TestX` | Corre só testes cujo nome dá match em `TestX` | `cargo test x` |
| `go test ./... -race` | Testes com detetor de race conditions ligado | não há flag direta, usarias ferramentas externas |
| `go test ./... -cover` | Mostra cobertura de testes | `cargo tarpaulin`/`cargo llvm-cov` (não é built-in) |
| `go vet ./...` | Análise estática — apanha erros comuns (ex: `Printf` com args errados) | `cargo check` (parcialmente) |
| `gofmt -l .` / `go fmt ./...` | Formata código Go (só há UMA forma "certa" — não há debate de estilo) | `cargo fmt` |
| `go doc <pkg>` | Mostra documentação de um pacote no terminal | `cargo doc` (mas sem gerar HTML) |
| `go list -m all` | Lista todas as dependências (diretas + transitivas) resolvidas | `cargo tree` |
| `go clean -modcache` | Limpa a cache local de módulos descarregados | `cargo clean` (parcial) |

Variáveis de ambiente importantes: `GOPATH` (onde ficam módulos descarregados e binários instalados, hoje raramente precisas de mexer), `GOPROXY` (proxy de módulos, default `proxy.golang.org`), `GOOS`/`GOARCH` (cross-compilation — ex: `GOOS=linux GOARCH=amd64 go build` gera um binário Linux a partir do teu Mac, sem toolchain extra, ao contrário do Rust onde precisas de instalar o target com `rustup target add`).

## Sintaxe — o essencial

### Variáveis e tipos
```go
var x int = 10       // declaração explícita
y := 20               // short declaration, tipo inferido (só dentro de funções)
const Pi = 3.14        // constante, resolvida em compile-time

var s string           // zero value: "" (Go não tem null/undefined — todo tipo tem "zero value")
var n int               // zero value: 0
var b bool              // zero value: false
var p *int              // zero value: nil (só ponteiros, slices, maps, chans, funcs, interfaces podem ser nil)
```
Não há `Option<T>` — o `nil` existe mas só para tipos "de referência" (ponteiros, slices, maps, channels, funcs, interfaces). Não podes ter `nil` num `int` ou `string` diretamente.

### Structs vs Rust structs/traits
```go
type URL struct {
    ID        uint
    ShortCode string
    Original  string
}

// "métodos" associam-se a um tipo via receiver, fora da struct:
func (u URL) IsExpired() bool { ... }     // receiver por valor (copia)
func (u *URL) IncrementVisits() { ... }    // receiver por ponteiro (muta o original)
```
Não há `impl Trait for Struct` bloco explícito — em Go qualquer tipo que tenha os métodos certos satisfaz uma interface automaticamente (**structural typing / duck typing**, resolvido em compile-time). Não declaras "URL implements Repository" em lado nenhum.

### Interfaces (equivalente a traits)
```go
type URLRepository interface {
    Save(u URL) error
    FindByCode(code string) (URL, error)
}
```
Qualquer struct que tenha esses métodos com essas assinaturas satisfaz a interface, sem `impl` explícito. Isto é o oposto do Rust, onde tens de declarar `impl Trait for Type` mesmo que os métodos já existam.

### Error handling (o maior contraste com Rust)
```go
func FindByCode(code string) (URL, error) {
    if code == "" {
        return URL{}, errors.New("código vazio")
    }
    return url, nil
}

// no caller:
u, err := FindByCode("abc")
if err != nil {
    return err   // ou tratar
}
```
Não há `Result<T, E>` nem `?`. Toda a função que pode falhar devolve `(T, error)` — dois valores — e **tu és responsável por verificar `err != nil` manualmente, sempre**. O compilador não te obriga a tratar o erro (ao contrário do `Result` do Rust, que é `#[must_use]`) — é possível ignorar um erro silenciosamente (`v, _ := f()`), o que é considerado má prática mas é legal. `errors.Is`/`errors.As` fazem o papel do pattern-matching em variantes de erro que terias com enums no Rust.

### Panics (equivalente a `panic!`)
```go
panic("algo correu muito mal")
```
Existe, mas é para erros verdadeiramente irrecuperáveis (bugs, invariantes quebradas) — não para erros esperados de negócio, que usam sempre `error`. `recover()` (só dentro de `defer`) apanha panics, papel parecido ao `catch_unwind` do Rust, mas usado sobretudo em frameworks (ex: middleware do gin) para não derrubar o servidor todo por um handler que rebentou.

### Ponteiros (sem borrow checker)
```go
x := 5
p := &x       // p é *int, aponta para x
*p = 10        // deref e escreve — x passa a ser 10
```
Go tem ponteiros mas **não tem ownership/borrow checker** — podes ter tantas referências quantas quiseres, ao mesmo objeto, ao mesmo tempo, mutáveis ou não, sem o compilador reclamar. Quem evita data races é a tua disciplina (ou `-race` a apanhar em runtime) e o **garbage collector**, que trata da memória — não precisas de `Drop`/`lifetimes`.

### Slices e maps (Vec e HashMap)
```go
s := []int{1, 2, 3}          // slice, ~Vec<i32>
s = append(s, 4)              // append devolve um slice novo (pode realocar)

m := map[string]int{"a": 1}   // ~HashMap<String, i32>
v, ok := m["a"]                // "comma ok idiom": ok diz se a key existia
delete(m, "a")
```
`ok` no acesso a maps é o equivalente funcional ao `Option` que o `.get()` do Rust devolveria — mas aqui são dois valores de retorno normais, não um tipo `Option<T>` empacotado.

### Goroutines e channels (vs async/tokio)
```go
go func() {
    fmt.Println("corre em paralelo")
}()

ch := make(chan int)      // channel não bufferizado
go func() { ch <- 42 }()   // envia
v := <-ch                   // recebe (bloqueia até haver valor)

select {
case v := <-ch1:
    ...
case v := <-ch2:
    ...
case <-time.After(time.Second):
    ...
}
```
`go func(){}()` arranca uma goroutine — gerida pelo **runtime do Go** (M:N scheduler, não é 1:1 com threads do SO), muitíssimo mais leve que uma thread (KBs de stack, cresce dinamicamente). Não há `async`/`await`/`Future` — código "assíncrono" em Go escreve-se de forma **síncrona/bloqueante** dentro da goroutine; é o scheduler que faz o multiplexing por baixo. `channel` é o mecanismo idiomático de comunicação entre goroutines ("Don't communicate by sharing memory; share memory by communicating" é o mote oficial do Go) — papel semelhante aos channels do `tokio::sync::mpsc`, mas built-in na linguagem com sintaxe própria (`<-`).

### Defer (RAII manual)
```go
f, err := os.Open("file.txt")
if err != nil { return err }
defer f.Close()   // corre no fim da função, aconteça o que acontecer (return normal, panic, etc.)
```
Não há `Drop` automático. `defer` é como agendares manualmente "corre isto quando esta função sair" — usado para fechar ficheiros, unlocks de mutex, fechar conexões, etc.

### Context (cancelamento cooperativo)
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

select {
case <-ctx.Done():
    return ctx.Err()
case result := <-doWork(ctx):
    return result
}
```
`context.Context` propaga-se explicitamente como primeiro argumento por convenção (`func Foo(ctx context.Context, ...)`) e carrega cancelamento/timeout/deadlines através de chamadas encadeadas — o equivalente a cancelar uma `Future`/task no tokio, mas cooperativo e explícito (tens de verificar `ctx.Done()` tu mesmo nos pontos certos).

## Ecossistema — ferramentas comuns

- **gin** (`github.com/gin-gonic/gin`) — framework HTTP mais popular, o que estamos a usar.
- **GORM** (`gorm.io/gorm`) — ORM mais usado.
- **chi** (`github.com/go-chi/chi`) — alternativa mais "minimalista"/próxima da stdlib ao gin.
- **sqlc** — gera código Go type-safe a partir de SQL puro (alternativa ao ORM, mais popular entre quem prefere SQL explícito).
- **testify** (`github.com/stretchr/testify`) — assertions mais ergonómicas para testes (`assert.Equal`, etc.) — a stdlib `testing` é minimalista de propósito.
- **golangci-lint** — agregador de linters (equivalente ao `clippy`, mas corre dezenas de linters diferentes por baixo).
- **air** — hot-reload para desenvolvimento local (recompila e reinicia ao gravar ficheiros), útil já que Go não tem "cargo watch" nativo mas tem este equivalente na comunidade.
- **viper** — gestão de configuração (env vars, ficheiros YAML/JSON, flags) mais robusta que `os.Getenv` manual.
- **cobra** — construir CLIs (usado pelo próprio `kubectl`, `hugo`, etc.).
- **golang.org/x/...** — extensões "quase-oficiais" da equipa Go que não entraram na stdlib (`x/time/rate` para rate limiting, `x/crypto`, `x/net`).
- **pkg.go.dev** — o equivalente ao `docs.rs`/`crates.io` combinados: documentação gerada automaticamente para qualquer módulo público.

## Convenções idiomáticas rápidas

- Nomes exportados (públicos) começam com **maiúscula** (`URL`, `FindByCode`); minúscula = privado ao pacote (`shortCode`). Não há keyword `pub` — é a capitalização que decide, ao nível de pacote (não de ficheiro).
- `gofmt` não é opcional na prática — quase toda a comunidade usa formatação automática, não há guerras de estilo tabs/espaços.
- Erros: mensagens em minúscula, sem pontuação final (`errors.New("code not found")`, não `"Code not found."`).
- Pacotes: nomes curtos, minúsculos, sem underscores (`shortcode`, não `short_code`).
