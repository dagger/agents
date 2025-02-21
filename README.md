# Melvin: a toy open-source programming agent

*WARNING: unfinished, experimental software*

This project is for educational purposes and should not be used for actual production-grade programming tasks.
There are several great "programming agent" products on the market, including [Devin](https://devin.ai), [Lovable](https://lovable.dev), [Replit Agent](https://replit.com), [V0](https://v0.dev) and many more.

## Architecture

Rather than a monolithic application, Melvin is a set of modular components that you can integrate into your application,
or use individually.

Each Melvin module has the following features:

- Runs in containers, for maximum portability and reproducibility
- Can be run from the command-line, or programmatically via an API
- Generated bindings for Go, Python, Typescript, PHP (experimental, Java (experimental), and Rust (experimental).
- End-to-end tracing of prompts, tool calls, and even low-level system operations. 100% of agent state changes are traced.
- Cross-language extensions. Add your own modules in any language.
- Platform-independent. No infrastructure lock-in! Runs on any hosting provider that can run containers.

These features are achieved by using the [Dagger Engine](https://dagger.io) as a runtime.
Dagger is open-source, and can be installed on any machine that can run Linux containers.
This includes Mac and Windows machines with a Docker-compatible tool installed.

Dagger is the only dependency for Melvin. No other tooling or programming environment is required:
the entire environment is containerized, for maximum portability.

## Modules

The following modules are currently implemented:

- [toy-workspace](./toy-workspace): a very, very simple development workspace for demo purposes
- [toy-programmer](./toy-programmer): a very, very simple programmer micro-agent for demo purposes
- [workspace](./workspace): a slightly more powerful workspace, with checkpoint and history features, and a configurable check function
- [reviewer](./reviewer): a code reviewer micro-agent
- [github](./github): a module for sending progress updates in a github issue
- [demo](./demo): a collection of demo functions tying the other modules together

## Community

We strongly recommend [joining our Community discord](https://discord.gg/KK3AfBP8Gw).
The Dagger community is very welcoming, and will be happy to answer your questions, discuss your use case ideas, and help you get started.

Do this now! It will make the rest of the experience more productive, and more fun.


## Initial setup

### 1. Install Dagger

*Note: the latest version is `0.17.0-llm.2`. It was released on Feb 21 2025. If you are running an older build, we recommend upgrading.*

Melvin's only dependency is Dagger - specifically a *development version* of Dagger which adds native support for LLM prompting and tool calling.

Once this feature is merged (current target is 0.17), Melvin will support with a stable release of Dagger.

Install the development version of LLM-enabled Dagger:

```console
curl -fsSL https://dl.dagger.io/dagger/install.sh | DAGGER_VERSION=0.17.0-llm.2 BIN_DIR=/usr/local/bin sh
```

You can adjust `BIN_DIR` to customize where the `dagger` CLI is installed.

Verify that your Dagger installation works:

```console
$ dagger core version
v0.17.0-llm.2
```

### 2.Configure LLM endpoints

Dagger uses your system's standard environment variables to route LLM requests. Currently these variables are supported:

- `OPENAI_API_KEY`
- `OPENAI_BASE_URL`
- `OPENAI_MODEL`
- `ANTHROPIC_API_KEY`
- `ANTHROPIC_BASE_URL`
- `ANTHROPIC_MODEL`

Dagger will look for these variables in your environment, or in a `.env` file in the current directory (`.env` files in parent directories are not yet supported).

#### Using with Ollama

You can use Ollama as a local LLM provider. Here's how to set it up:

1. Install Ollama from [ollama.ai](https://ollama.ai)

2. Start Ollama server with host binding:

```shell
OLLAMA_HOST="0.0.0.0:11434" ollama serve
```

3. Get your host machine's IP address:

On macOS / Linux (modern distributions):

```shell
ifconfig | grep "inet " | grep -v 127.0.0.1
```

On Linux (older distributions):

```shell
ip addr | grep "inet " | grep -v 127.0.0.1
```

This step is needed because our LLM type runs inside the engine and needs to reach your local Ollama service. While we're potentially exploring the implementatin of automatic tunneling, for now you'll need to use your machine's actual IP address instead of localhost to allow the containers to communicate with Ollama.

1. Configure the following environment variables (replace `YOUR_IP` with the IP address from step 3):

```plaintext
OPENAI_API_KEY="nothing"
OPENAI_BASE_URL=http://YOUR_IP:11434/v1/
OPENAI_MODEL=llama3.2
```

For example, if your IP is 192.168.64.1:
```plaintext
OPENAI_API_KEY="nothing"
OPENAI_BASE_URL=http://192.168.64.1:11434/v1/
OPENAI_MODEL=llama3.2
```

## Run Melvin from the command-line

To run Melvin from the command-line, use the `dagger` CLI to load one of its modules, and call its functions.

For example, to use the `toy-programmer` module:

```console
dagger shell -m ./toy-programmer
```

Then, run this command in the dagger shell:

```
.doc
```

This prints available functions. Let's call one:

```
go-program "develop a curl clone" | terminal
```

This calls the `go-program` function with a description of a program to write, then runs the `terminal` function on the returned container.

You can use tab-completion to explore other available functions.

### Integrate Melvin in your application

You can embed Dagger modules into your application.
Supported languages are Python, Typescript, Go, Java, PHP - with more language support under development.

1. Initialize a Dagger module at the root of your application.
This doesn't need to be the root of your git repository - Dagger is monorepo-ready.

```console
dagger init
```

2. Install the modules you wish to load

For example, to install the Melvin toy-workspace module:

```console
dagger install github.com/shykes/melvin/toy-workspace
```

3. Install a generated client in your project

*TODO: this feature is not yet merged in a stable version of Dagger*

This will configure Dagger to generate client bindings for the language of your choice.

For example, if your project is a Python application:

```console
dagger client install python
```

4. Re-generate clients

*TODO: this feature is not yet merged in a stable version of Dagger*

Any time you need to re-generate your client, run:

```console
dagger client generate
```
