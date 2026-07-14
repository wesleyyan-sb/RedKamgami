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
