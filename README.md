# onekube

Store your kubeconfig files within [1password](https://1password.com/) and pull them down on demand

## How it works

`onekube` uses the 1password CLI to download kubeconfig files tagged in your 1password vaults to be used locally.

## Getting started

### Installation

Clone the repostory and run:

```
go install .
```

Ensure your $PATH includes your go binaries

### Completions (optional)

Auto completions are available for several shells, run the below to list them
```
onekube completion -h
```

More instructions to add completions can be found for a specific shell like so
```
onekube completion zsh -h
```

### Using onekube

#### 1password vault

All kubeconfigs should be stored in 1password with the contents of the kubeconfig in a field called `config` and a tag `kubeconfig`.
The name/selector used by onekube will be the name of the item.

#### onekube CLI

To start with a clean slate, run
```
onekube init
```

View available configs with
```
onekube list
```

Select a config with
```
onekube set <my-config>
```

Your existing config will be backed up, restore it with either
```
onekube init
```
or
```
onekube clean
```

> [!WARNING]
> Any kubeconfig files set will be stored on your local machine until replaced with another, or `onekube clean` is run

## Dependencies

- The [1password CLI](https://developer.1password.com/docs/cli) must be installed and logged in

## Limitations / Known issues

- If you have multiple 1password accounts you must change between these with the 1password CLI for onekube to use them
- If you are not logged into 1password, the experience is poor
