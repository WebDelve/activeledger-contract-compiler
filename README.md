<img src="assets/webdelve-logo.png" alt="WebDelve"/><br/>
[WebDelve - Purpose Drive Software](https://webdelve.co)

*** 

<img src="https://github.com/activeledger/activeledger/blob/master/docs/assets/Asset-23.png" alt="Activeledger" width="250"/><br/>
This project is built to work with Activeledger<br/>
[Activeledger on GitHub]( https://github.com/activeledger/activeledger )<br/>
[Activeledger Website](https://activeledger.io)

***

# Contract Compiler

This software allows you to write Activeledger Smart Contracts in multiple files
and compile them into one for uploading to a network.

Rather than writing a complex contract which utilises multiple classes in a 
single file, you can break those classes out into separate files. This software
will look for local imports in the entry file provided by the user and merge
external imports and the classes in each local imported file into one output.

The external imports are handled such that
```typescript
import { Activity } from "@activeledger/activeledgercontracts";
```
from one file, and
```typescript
import { Standard } from "@activeledger/activeledgercontracts";
```
from another, become
```typescript
import { Activity, Standard } from "@activeledger/activeledgercontracts";
```
in the output.

If multiple files have
```typescript
import { Activity, Standard } from "@activeledger/activeledgercontracts";
```
that import line will only be included once in the output.

## Quick start

### Building

With Go installed run
```bash
make build
```
to build the software. This will create a build in `./bin`.

Move the build and the config file local to the folder containing your Smart
Contract.

### Running

Two CLI flags are available to be used.
```bash
# Provide an entry file
./comp -p <smartcontractfolder>/<entryfile>.ts

# Provide an output file path
./comp -o <output>.ts

# Provide both
./comp -p <smartcontractfolder>/<entryfile>.ts -o <output>.ts
```
Example:
```bash
./comp -p smartcontract/main.ts -o compiled.ts
```
This will look for a file called `main.ts` in the local directory `smartcontract`,
and output to the local file `compiled.ts`.

**Note:** `compiled.ts` will be created if it doesn't exist, and overwritten if
it does

## Contract files
Currently the software will look for local imports like the following
```typescript
import { Second } from "./second";
```
Where `Second` is a class exported from a file named `second.ts` in the same 
directory as the provided entry file.




