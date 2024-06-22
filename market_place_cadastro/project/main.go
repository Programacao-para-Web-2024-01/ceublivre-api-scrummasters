package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"

    _ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type User struct {
    Nome                   string `json:"nome"`
    Nascimento             string `json:"nascimento"`
    Contato                string `json:"contato"`
    Email                  string `json:"email"`
    Senha                  string `json:"senha"`
    ConfirmarSenha         string `json:"confirmar_senha"`
    Cep                    string `json:"cep"`
    Endereco               string `json:"endereco"`
    Numero                 string `json:"numero"`
    Complemento            string `json:"complemento"`
    Bairro                 string `json:"bairro"`
    Cidade                 string `json:"cidade"`
    Uf                     string `json:"uf"`
    PreferenciasComunicacao string `json:"preferencias_comunicacao"`
}

// Estrutura para representar um vendedor
type Vendedor struct {
    Nome                 string `json:"nome"`
    Email                string `json:"email"`
    Senha                string `json:"senha"`
    ConfirmarSenha       string `json:"confirmar_senha"`
    Banco                string `json:"banco"`
    Agencia              string `json:"agencia"`
    Conta                string `json:"conta"`
    InformacoesFiscais   string `json:"informacoes_fiscais"`
    InformacoesBancarias string `json:"informacoes_bancarias"`
}

func main() {
    var err error
    db, err = sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/market_place")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    err = db.Ping()
    if err != nil {
        log.Fatal(err)
    }

    http.HandleFunc("/registerUser", registerUserHandler)
    http.HandleFunc("/registerVendedor", registerVendedorHandler)

    fmt.Println("Server started at :8081")
    log.Fatal(http.ListenAndServe(":8081", nil))
}

func registerUserHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        var user User
        err := json.NewDecoder(r.Body).Decode(&user)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        if user.Senha != user.ConfirmarSenha {
            http.Error(w, "As senhas não coincidem", http.StatusBadRequest)
            return
        }

        stmt, err := db.Prepare("INSERT INTO usuarios(nome, nascimento, contato, email, senha, cep, endereco, numero, complemento, bairro, cidade, uf, preferencias_comunicacao) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        defer stmt.Close()

        _, err = stmt.Exec(user.Nome, user.Nascimento, user.Contato, user.Email, user.Senha, user.Cep, user.Endereco, user.Numero, user.Complemento, user.Bairro, user.Cidade, user.Uf, user.PreferenciasComunicacao)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Usuário cadastrado com sucesso"))
    } else {
        http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
    }
}

func registerVendedorHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        var vendedor Vendedor
        err := json.NewDecoder(r.Body).Decode(&vendedor)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        if vendedor.Senha != vendedor.ConfirmarSenha {
            http.Error(w, "As senhas não coincidem", http.StatusBadRequest)
            return
        }

        stmt, err := db.Prepare("INSERT INTO vendedores(nome, email, senha, banco, agencia, conta, informacoes_fiscais, informacoes_bancarias) VALUES(?, ?, ?, ?, ?, ?, ?, ?)")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        defer stmt.Close()

        _, err = stmt.Exec(vendedor.Nome, vendedor.Email, vendedor.Senha, vendedor.Banco, vendedor.Agencia, vendedor.Conta, vendedor.InformacoesFiscais, vendedor.InformacoesBancarias)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Vendedor cadastrado com sucesso"))
    } else {
        http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
    }
}
