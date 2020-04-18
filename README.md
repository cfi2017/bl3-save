[![GoDoc](https://godoc.org/github.com/cfi2017/bl3-save?status.svg)](https://godoc.org/github.com/cfi2017/bl3-save)
[![Go Report Card](https://goreportcard.com/badge/github.com/cfi2017/bl3-save)](https://goreportcard.com/report/github.com/cfi2017/bl3-save)

# bl3-save

Command line utility for modifying borderlands 3 character and profile saves.

## Getting started

To get the utility, go to the releases page or run `go get github.com/cfi2017/bl3-save`.

**If you're using Steam, make sure to disable cloud synchronisation or Steam will overwrite your changes with your cloud save.**

**Make sure to make a backup before using this tool. The author takes no responsibility for any loss of progress as a result of this.**

## Using
- Extract the files from the release into your Borderlands 3 Save directory. 
- Run bl3-save.exe.

## Backups
This tool automatically creates backups of your save files whenever you save them. 
You can find these backups in the backup directory of your save files.

## Credits

Credits go to these amazing people for all various sorts of reasons:
- Gibbed & Apocalyptech for their work on decrypting and decoding save files and items
- Cu3PO42 for expanding on the roadmap and their input regarding various features of the client.
- aprizm for collaboration for conversion between items in various formats and base64 encoded json

### Bug hunters
- Zydiz
