
# KEC Manual

## Introduction
KEC *(short for K-mer exclusion by crossreference)* was designed to search for unique DNA / RNA / amino acid sequences in large datasets. The original aim was to find unique sites to design (PCR / LAMP etc.) primers for specific detection of bacteria. It takes **target** and **non-target** genomes and by crossreferencing of the respective K-mers it reconstructs sequences that are unique for **target** genome(s).

## Principle of operation

## Installation
KEC does not require installation. Binary executables for Windows, Linux and macOS can be downloaded from Releases section. After download and extraction the software works in any directory. During download and first run the software can be marked by antivirus as harmul. However, if downloaded from this repository, the program only works as stated, without any malicious activity or data collection. Users may inspect and compile from source code, if security is a concern.

## Command-line parameters
KEC has two modes of operation - **exclude** and **include**. Each mode has its own set of parameters accessible by `-h` or `--help` parameter (e.g. `kec.exe exclude -h`).

### exclude mode
- `-t` - Add target sequence(s). Fasta formatted file or whole directory is allowed.
- `-n` - Add nontarget sequence(s). Fasta formatted file or whole directory is allowed.
- `-o` - Output path to store resulting fasta formatted file.
- `-k` - K-mer size. Explained in separate chapter. *Default = 12*.
- `-r` - Also exclude reverse complements of the sequences. Takes more (approx. 2 - 3x) time to finish. *Default = false*.
- `--min` - Minimum size of resulting sequence. *Default = 13*.
- `--max` - Maximum size of resulting sequence (0 = unlimited). *Default = 0*.
- `--help`, `-h` - Show help message.

**Example:**

`kec.exe exclude -t c:\seqs\target -n c:\seqs\nontarget -o c:\seqs\results\test.fasta -k 12 --min 200`

Will search for unique sequences in fasta formatted files located in `c:\seqs\target` by comparing them to non-target sequences in `c:\seqs\nontarget`. K-mer size will be 12 and minimal length of recovered sequences will be 200.

### include mode
- `-m` - Add master sequence(s). Fasta formatted file or whole directory is allowed.
- `-p` - Add pool sequence(s) to include in consensus sequence. Fasta formatted file or whole directory is allowed.
- `-o` - Output path to store resulting fasta formatted file.
- `-k` - K-mer size. Explained in separate chapter. *Default = 12*.
- `--min` - Minimum size of resulting sequence. *Default = 13*.
- `--max` - Maximum size of resulting sequence (0 = unlimited). *Default = 0*.
- `--help`, `-h` - Show help message.

**Example:**

`kec.exe include -m c:\seqs\master\master.fasta -n c:\seqs\nontarget -o c:\seqs\results\test.fasta -k 12 --min 200`

Will search for common sequences in fasta formatted file `c:\seqs\master\master.fasta` by comparing them to pool sequences located in directory `c:\seqs\pool`. K-mer size will be 12 and minimal length of recovered sequences will be 200.

## K selection

## Input and output data
KEC currently works only with fasta formated files on both input and output. To work with other formats (e.g. fastq) the file has to be converted to fasta first by any available tool. KEC also accepts whole directory as an input, where files with extensions `.fasta` `.fna`, `.ffn`, `.faa` and `.frn` are used. Direct support for other formats will be added in future versions.


## Memory tests

## Speed tests

## Citation
Under review