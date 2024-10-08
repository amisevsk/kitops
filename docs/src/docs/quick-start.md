<script setup>
import vGaTrack from '@theme/directives/ga'
</script>

# KitOps Quick Start

In this guide, we'll use ModelKits and the kit CLI to easily:
* Package up a model, notebook, and datasets into a single ModelKit you can share through your existing tools
* Push the ModelKit package to a public or private registry
* Grab only the assets you need from the ModelKit for testing, integration, local running, or deployment
* Run an LLM locally to speed app integration, testing, or experimentation

## Before we start...

1. Make sure you've got the [Kit CLI setup](./cli/installation.md).
2. Create and navigate to a new folder on your filesystem - we suggest calling it `KitStart` but any name works.

## Learning to use the CLI

### 1/ Check your CLI Version

Check that the Kit CLI is properly installed by using the [version command](./cli/cli-reference.md#kit-version).

```sh
kit version
```

You'll see information about the version of Kit you're running. If you get an error check to make sure you have [Kit installed](./cli/installation.md) and in your path.

### 2/ Login to Your Registry

You can use the [login command](./cli/cli-reference.md#kit-login) to authenticate with any OCI v1.1-compatible container registry. Here we'll pull ModelKits from [Jozu Hub](https://jozu.ml/discover), but we'll use **GitHub Registry** to push (the Jozu Hub only supports pull, although push is being worked on now).

```sh
kit login ghcr.io
```

After entering your username and password, you'll see `Log in successful`. If you get an error it may be that you need an HTTP vs HTTPS (default) connection. Try the login command again but with `--plain-http`.

### 3/ Get a Sample ModelKit

Let's use the [unpack command](./cli/cli-reference.md#kit-unpack) to pull a [sample ModelKit](./modelkit/premade-modelkits.md) to our machine that we can play with. In this case we'll unpack the whole thing, but one of the great things about Kit is that you can also selectively unpack only the artifacts you need: just the model, the model and dataset, the code, the configuration...whatever you want. Check out the `unpack` [command reference](./cli/cli-reference.md#kit-unpack) for details.

You can grab <a href="https://jozu.ml/discover"
  v-ga-track="{
    category: 'link',
    label: 'grab any of the ModelKits',
    location: 'docs/quick-start'
  }">any of the ModelKits</a> from the Jozu Hub, but we've chosen a fine-tuned model based on Llama3.

```sh
kit unpack jozu.ml/jozu/fine-tuning:tuned
```

You'll see a set of messages as Kit unpacks the configuration, code, datasets, and serialized model. Now list the directory contents:

```sh
ls
```

You'll see:
* A Llama3 model
* A LoRA adapter
* A training dataset
* A README file
* A Kitfile

The [Kitfile](./kitfile/kf-overview.md) is the manifest for our ModelKit, the serialized model, and a set of files or directories including the adapter, dataset, and docs. Every ModelKit has a Kitfile and you can use the info and inspect commands to view them from the CLI (there's more on this in our [Next Steps](next-steps.md) doc).

### 4/ Check the Local Repository

Use the [list command](./cli/cli-reference.md#kit-list) to check what's in our local repository.

```sh
kit list
```

You'll see the column headings for an empty table with things like `REPOSITORY`, `TAG`, etc...

### 5/ Pack the ModelKit

Since our repository is empty we'll need use the [pack command](./cli/cli-reference.md#kit-pack) to create our ModelKit. The ModelKit in your local registry will need to be named the same as your remote registry. So the command will look like: `kit pack . -t [your registry address]/[your repository name]/mymodelkit:latest`

In my case I am pushing to the `jozubrad` repository:

```sh
kit pack . -t ghcr.io/jozubrad/mymodelkit:latest
```

You'll see a set of `Saved ...` messages as each piece of the ModelKit is saved to the local repository.

Checking your local registry again you should see an entry:

```sh
kit list
```

The new entry will be named based on whatever you used in your pack command.

### 6/ (Optional) Remove a ModelKit from a Local Repository

Let's pretend that the `pack` command we ran in the previous step contained a typo in the ModelKit's repository name causing the word "model" to be entered as "modle". The output from the `kit list` command would display the ModelKit as:

```sh
ghcr.io/jozubrad/mymodlekit:latest
```

To correct this, we would `remove` the misspelled ModelKit from our local repository using the [remove command](./cli/cli-reference.md#kit-remove), being sure to provide reference the ModelKit using its mispelled name:

```sh
kit remove ghcr.io/jozubrad/mymodlekit:latest
```

Next, we would repeat the `kit pack` command in the previous step, being sure to provide the correct repository name for our ModelKit.

### 7/ Push the ModelKit to a Remote Repository

The [push command](./cli/cli-reference.md#kit-push) will copy the newly built ModelKit from your local repository to the remote repository you logged into earlier. The naming of your ModelKit will need to be the same as what you see in your `kit list` command (REPOSITORY:TAG). You can even copy and paste it. In my case it looks like:

<!-- replace with Jozu Hub once private repos are ready -->

```sh
kit push ghcr.io/jozubrad/mymodelkit:latest
```

### Congratulations

You've learned how to unpack a ModelKit, pack one up, push it, and run an LLM locally. Anyone with access to your remote repository can now pull your new ModelKit and start playing with your model using the `kit pull` or `kit unpack` commands.

If you'd like to learn more about using Kit, try our [Next Steps with Kit](./next-steps.md) document that covers:
* Signing your ModeKit
* Making your own Kitfile
* The power of `unpack`
* Tagging ModelKits
* Keeping your registry tidy

Thanks for taking some time to play with Kit. We'd love to hear what you think. Feel free to drop us an [issue in our GitHub repository](https://github.com/jozu-ai/kitops/issues) or join [our Discord server](https://discord.gg/3eDb4yAN).
