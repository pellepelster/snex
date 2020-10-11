![snex](https://github.com/pellepelster/snex/workflows/snex/badge.svg)

# SNippet EXtractor (snex)
This utility helps you with the task of keeping code samples inside of documentation in sync with real world code from your sources. The issues this solves is that source examples in documentation tends to get outdated very quick. By pulling the source directly from a working project you can make sure the source examples used in you docs are always up to date.

## Downloads
* Linux (AMD64) [snex_linux_amd64](https://github.com/pellepelster/snex/releases/latest/download/snex_linux_amd64)
* Windows (AMD64) [snex_windows_amd64](https://github.com/pellepelster/snex/releases/latest/download/snex_windows_amd64)
* Darwin (AMD64) [snex_darwin_amd64](https://github.com/pellepelster/snex/releases/latest/download/snex_darwin_amd64)

## How it works
`snex` consumes a list of files line by line, and looks for lines containing snippet start- and end-markers. Then it crawls through all documentation files (e.g. a list of markdown files in a folder) and replaces all references to code snippets by the actual code. To keep things simple and language agnostic it does not care for comment markers (which differ between languages) and just looks for snippet start- and end-markers.

So assuming you are tasked with writing some documentation and want to ensure it always contains up to date and correct examples. Just mark the place where the snippet should be inserted:
```markdown
# Documentation

How to assign names to variables:

<!--- snippet:snippet1 -->
<!--- /snippet:snippet1 -->

this is how you use strings in Java.

```

then in your code/unit tests mark the code accordingly

```java
public class ClassContainingSnippets {
    public void doSomeImportantStuff() {
        // snippet:snippet1
        var firstName = "Jens";
        var lastName = "Mander";
        // /snippet:snippet1
    }
}
```

when you let `snex` run over your documentation it will then replace the snippet markers with the real code:

```markdown
# Documentation

How to assign names to variables:

<!--- snippet:snippet1 -->
        var firstName = "Jens";
        var lastName = "Mander";
<!--- /snippet:snippet1 -->

this is how you use strings in Java.
```

`snex` will keep the original markers to ensure it can be re-run anytime on the documentation sources. 

## Usage

### Folder based

```shell script
$ snex -source ./documentation -target ./output -snippets ./src
```

where **source** is the directory containing the files to transform, **target** is the output directory containing the transformed files and **snippets** designates the directory containing the code samples

### Single file based

```shell script
$ snex -source ./src/README.md -snippets ./src
```

if a single file is specified as **source** and no target is given the file is transformed in place, using the snippets found in **snippets**

### Template support
To support various uses cases, a template can be given that is used to render the snippets, for example to render all snippets inside a markdown code block use this template:

```shell script
snex --source ./POST.md  -snippets ./  -template '```{{.Content}}```'
```

This will result in the following output:

```shell script
# Documentation

How to assign names to variables:

<!--- snippet:snippet1 -->
    ```
        var firstName = "Jens";
        var lastName = "Mander";
    ```
<!--- /snippet:snippet1 -->

this is how you use strings in Java.
```

for a more complex templates you can also use a file, see this example that was used to generate [this](https://github.com/pellepelster/pelle.io/blob/master/site/content/posts/ca-secured-ssh-connections.md) hugo post:

```shell script
snex --source ./POST.md  -snippets ./  -template-file hugo.template
```

**hugo.template**
```markdown
{{ if .IsFullFile }}
{{`{{< github repository="pellepelster/vault-ssh-ca"`}} file="{{.Filename}}"  >}}{{.Filename}}{{`{{< /github >}}`}}
{{ else }}
{{`{{< github repository="pellepelster/vault-ssh-ca"`}} file="{{.Filename}}#L{{.Start}}-L{{.End}}"  >}}{{.Filename}}{{`{{< /github >}}`}}
{{ end }}
{{`{{< highlight go "" >}}`}}
{{.Content}}
{{`{{< / highlight >}}`}}
```
