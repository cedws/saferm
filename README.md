# saferm
Utility for intercepting C standard library calls to delete files, instead moving them to the `Trash` directory. This is achieved via `LD_PRELOAD`. This can prevent an accidental `rm` or software bug from becoming a disaster.

**This is just an experiment for now. Use at your own risk.**

## Limitations
- Can only prevent dynamically linked executables from truly deleting your files
- Only supports Linux (macOS support might be possible but difficult, due to runtime hardening)

## Usage
- Run `make`
- Configure your shell to export `LD_PRELOAD` with the path to the built shared object:

    - Bash: add to `~/.profile`:
        ```
        export LD_PRELOAD=full/path/to/saferm.so
        ```

    - Fish: add to `~/.config/fish/config.fish`:
        ```
        set -x LD_PRELOAD full/path/to/saferm.so
        ```
- Reload your shell configuration with `source`
