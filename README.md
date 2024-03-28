<h3 align=center><b>untree</b></h3>

<b>untree</b> finds and flattens tree-like text, making it easily searchable. for example prefixing all expressions a codebase with their parent function definitions, all HTML tags with their parent tags, or JSON objects their parent keys names.

all that's required is indentation with spaces or tabs

<i>(like <a href="https://github.com/tomnomnom/gron">gron</a>, but generalised on indentation)</i>

---

### installation

```
    $ go install go.senan.xyz/untree@latest
```

### usage

```
    $ cmd | untree
    $ untree [FILE ...]
```

the output format is `<prefix>\t<original line>`. suitable for grepping, piping, column selecting. eg "all Go handler functions which log under a certain condition" could be

```bash
untree "$(git ls-files "*.go")" \
    | grep "func.*Handle.*if.*err :=.*slog.Info" \
    | awk '{ print $2 }'`
```

### examples

<p float="left" align=center>
  <img src="https://github.com/sentriz/untree/assets/6832539/5af5ed5f-e0a0-4f10-a633-2f5a4adfc5cb" />
  <img src="https://github.com/sentriz/untree/assets/6832539/a15c8adf-6319-4ce1-9d17-f88d441c335e" />
</p>
