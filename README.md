# SNippet EXtractor (snex)
This little utility helps you with the task of keeping code samples inside of documentation in sync with real world code that is proven to work.

## How it works
`snex` consumes files line by line, and looks for lines containing snippet start- and end-markers. To keep things simple and language agnostic it does not care for comment markers depending on the language.


So assuming you are tasked with writing some documentation and want to ensure it always contains up to date and correct examples. Just mark the place where the snippet should be inserted:
```markdown
# Documentation

How to assign names to variables:

\<!--- snippet:snippet1 -->
\<!--- /snippet:snippet1 -->

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

when you let `snex` run over you documentation it will replace the snippet markers with the real code:

```markdown
# Documentation

How to assign names to variables:

\<!--- snippet:snippet1 -->
        var firstName = "Jens";
        var lastName = "Mander";
\<!--- /snippet:snippet1 -->

this is how you use strings in Java.

```

`snex` will keep the markers to ensure it can be re-run anytime on the documentation. 

## Usage

```shell script
$ snex -source ./documentation -target ./output -snippets ./src
```

where **source** is the directory containing the files to transform, **target** is the output directory containing the transformed files and **snippets** designates the directory containing the code samples