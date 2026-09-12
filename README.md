<p align="center">
  <a href="https://github.com/evanrock520-ciencias/Gabby/actions/workflows/ci.yml">
    <img src="https://github.com/evanrock520-ciencias/Gabby/actions/workflows/ci.yml/badge.svg" alt="CI State">
  </a>
  <img src="https://img.shields.io/badge/-Rust-3776AB?style=flat&logo=rust" alt="Rust">
  <img src="https://img.shields.io/badge/Go-00ADD8?logo=go&logoColor=white" alt="Go">
</p>

# Gabby

**Gabby** es un proyecto de mensajería personal por terminal con arquitectura cliente-servidor. Es desarrollado como primer proyecto del curso "**Modelado y Programación**" de la **Facultad de Ciencias, UNAM**.

## Cliente

El cliente está desarrollado con el lenguaje de programación **Go** utilizando una gran cantidad de bibliotecas proporcionadas por **Charm**. Cuenta con una interfaz de
terminal basada en aplicaciones de terminal modernas como [Lazygit](https://github.com/jesseduffield/lazygit) y [Lazydocker](https://github.com/jesseduffield/lazydocker) 
organizada siguiendo la filosofía de la **Arquitectura ELM**.



### Dependencias

- [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- [Lipgloss](https://github.com/charmbracelet/lipgloss)
- [Bubbles](https://github.com/charmbracelet/bubbles)

## Servidor

El servidor será desarrollado con el lenguaje de programación **Rust**. La elección 
radica en la gran velocidad y seguridad que ofrece de fábrica el lenguaje y en la alta
capacidad para sistemas concurrentes de la biblioteca **Tokio**

### Dependencias

- [Tokio](https://github.com/tokio-rs/tokio)
- [Serde](https://github.com/serde-rs/serde)

## Cómo comenzar

Para compilar el proyecto necesitamos tener instaladas las siguientes dependencias en
nuestro sistema.

### Dependencias

- [Go](https://github.com/golang/go)
- [Rust](https://github.com/rust-lang/rust)
- [Just](https://github.com/casey/just)

### Compilar

Para compilar el proyecto, desde la raíz del repositorio:

```bash
just build
```

Si solamente queremos compilar alguno de los 2 binarios:

`server`

```bash
just build-server
```

`client`

```bash
just build-client
```

### Pruebas

Para correr las pruebas, desde la raíz del repositorio:

```bash
just test
```
Si solamente queremos correr pruebas alguno de los 2 módulos:

`server`

```bash
just test-server
```

`client`

```bash
just test-client
```

### Correr

`server`

```bash
just run-server
```

`client`

```bash
just run-client
```

## Licencia

Este proyecto se encuentra bajo la licencia [GPL-3](LICENSE).