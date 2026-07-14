# RedKamgami

<p align="center">
  <h1 align="center">🔴 RedKamgami</h1>
  <p align="center">
    A large-scale password wordlist
  </p>
</p>

---

## Disclaimer

**RedKamgami is intended exclusively for:**

* Security research
* Authorized penetration testing
* Password auditing
* Academic research
* Defensive cybersecurity activities

The author does **not** encourage or endorse unauthorized access to computer systems, networks, or accounts. Users are solely responsible for ensuring that their use of this project complies with all applicable laws and regulations.

---

# Installation & Usage

RedKamgami supports **Windows**, **Linux**, and **macOS**.

## 1. Download the Launcher

Go to the project's **Releases** page and download the launcher for your operating system.

---

## 2. Make the Launcher Executable (Linux/macOS)

If you're using Linux or macOS, grant execute permission:

```bash
chmod +x RedKamgami-linux-amd64
```

or

```bash
chmod +x RedKamgami-macos-amd64
```

---

## 3. Launch RedKamgami

Run the launcher.

**Windows**

```powershell
.\RedKamgami-windows-amd64.exe
```

**Linux**

```bash
./RedKamgami-linux-amd64
```

**macOS**

```bash
./RedKamgami-macos-amd64
```

---

## 4. Download the Wordlist

From the launcher menu, choose the option to download the latest wordlist.

The launcher will automatically:

* Select an available mirror
* Verify the download
* Display download progress
* Save the compressed archive (`.zst`)

---

## 5. Place the Downloaded File

After the download is complete, ensure the downloaded `.zst` file is located in the **same directory as the launcher**.

Example:

```text
RedKamgami/
├── RedKamgami-linux-amd64
└── RedKamgami-2026.07.zst
```

---

## 6. Extract the Wordlist

Run the launcher again and choose the **Extract Wordlist** option.

The launcher will:

* Verify the archive integrity
* Decompress the `.zst` archive
* Produce the final wordlist
* Optionally remove the compressed archive after successful extraction

Result:

```text
RedKamgami/
├── RedKamgami-linux-amd64
└── RedKamgami.txt
```

---

## Updating

Whenever a new release is available:

1. Launch RedKamgami.
2. The launcher checks for updates automatically.
3. Download the latest compressed wordlist.
4. Replace the old archive (if applicable).
5. Use the launcher to extract the updated version.

No manual decompression tools are required.

---

## Supported Platforms

| Operating System | Supported |
| ---------------- | --------- |
| Windows          | ✅         |
| Linux            | ✅         |
| macOS            | ✅         |

The launcher is distributed as a native executable for each supported platform and requires no additional runtime or dependencies.


# Features

* Massive password collection
* Duplicate-free dataset
* High-performance password generator (Rust)
* Automatic launcher/updater (Go)
* Multiple download mirrors
* SHA-256 integrity verification
* Compressed distribution using Zstandard (.zst)
* Open source
* Cross-platform launcher

---

# Components

## Password Dataset

The primary dataset consists of a curated collection of passwords assembled for research purposes.

The build pipeline includes:

* Dataset merging
* Duplicate removal
* Validation
* Statistics generation
* Compression

---

## Password Generator (Rust)

A high-performance CLI application capable of generating millions of passwords with configurable options.

Supported modes include:

* Numeric only
* Alphanumeric
* Alphanumeric + symbols
* Fully customizable character sets

Designed for:

* Speed
* Low memory usage
* Large-scale generation

---

## Launcher (Go)

The launcher automatically handles:

* Version checking
* Mirror selection
* Download management
* Resume support (when available)
* SHA-256 verification
* Automatic extraction

Supported mirrors may include:

* Google Drive
* MEGA
* MediaFire
* OneDrive
* TeraBox

---

# Compression

The distributed wordlists use **Zstandard (.zst)** to reduce download size while maintaining extremely fast decompression.

---

# Integrity Verification

Every release includes a SHA-256 checksum.

The launcher verifies downloaded files automatically before extraction.

---

# Performance

The password generator is designed for:

* Streaming output
* Buffered file writing
* Minimal memory allocation
* Multi-core CPUs (where beneficial)

---

# Intended Use

RedKamgami is designed for:

* Password auditing
* Security research
* Academic projects
* Digital forensics
* Authorized penetration testing
* Defensive security exercises

---

# License

This project is released under the **Creative Commons Zero 1.0 Universal (CC0 1.0)**.

You are free to:

* Use
* Modify
* Redistribute
* Incorporate into other projects

without restriction.

For full license details, see the LICENSE file.

---

# Repository Policy

This repository is maintained as a personal project.

## Pull Requests

Pull requests are **not accepted**.

## Issues

Issue reports are **disabled**.

Suggestions, feature requests, and bug reports are not tracked through this repository.

---

# Acknowledgements

This project is inspired by the broader cybersecurity research community and by publicly available research on password security, password hygiene, and defensive security testing.

---

# Support

If you find RedKamgami useful, consider giving the repository a ⭐.

---

<p align="center">
Made with Rust, Go, and far too many passwords.
</p>
