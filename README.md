# Ethereum Address Generator (genaddr.exe)

A fast and flexible Ethereum address generator with pattern matching capabilities. Generate Ethereum addresses that match specific patterns at the beginning, end, or anywhere in the address.

## Features

- Advanced pattern-based address generation
- Support for multiple patterns
- Complex pattern matching with multiple wildcards
- Special character support for digits and letters
- Multi-threaded processing
- Real-time statistics
- Continuous search mode
- Save results to file
- Case-insensitive pattern matching

## Usage

```bash
genaddr.exe [options]
```

### Options

- `-pattern string` (required)
  Pattern to match. Supports multiple patterns separated by commas.
  Examples:
  - Simple patterns:
    - `123*` : address starts with "123"
    - `*123` : address ends with "123"
    - `*123*` : address contains "123"
  - Complex patterns:
    - `123*321` : address starts with "123" and ends with "321"
    - `1*2*3*4` : address contains "1", "2", "3", "4" in sequence
    - `*dead*beef*` : address contains "dead" followed by "beef"
  - Special characters:
    - `#` : matches any digit (0-9)
    - `@` : matches any letter (a-f)
    Examples:
    - `###*` : starts with any three digits
    - `@@@*` : starts with any three letters
    - `#@#@*` : starts with alternating digit and letter
    - `*###` : ends with any three digits
    - `##@@##*` : starts with two digits, then two letters, then two digits
  - Multiple patterns:
    - `123*,*456,*789*` : matches any of these patterns
    - `dead*,*cafe*,babe*` : matches addresses starting with "dead" OR containing "cafe" OR starting with "babe"
    - `###*,@@@*` : matches addresses starting with either three digits OR three letters
- `-workers int` (default: 4)
  - Number of worker goroutines
- `-continue`
  - Continue searching after finding a match
- `-output string`
  - Save found addresses to file
- `-help`
  - Show help message

### Examples

1. Find an address with specific start and end:
```bash
genaddr.exe -pattern "dead*beef"
```

2. Find addresses matching any of multiple patterns:
```bash
genaddr.exe -pattern "cafe*,*dead*,*babe" -workers 8 -continue -output results.txt
```

3. Find an address starting with three digits followed by a letter:
```bash
genaddr.exe -pattern "###@*" -workers 4
```

4. Find an address with alternating digits and letters at start:
```bash
genaddr.exe -pattern "#@#@#@*" -workers 8
```

## Output Format

For each found address, the program outputs:
```
Address: 0x123...
Private Key: 0x456...
```

## Performance Tips

1. Adjust the number of workers based on your CPU cores:
   - For 4-core CPU: use `-workers 4`
   - For 8-core CPU: use `-workers 8`
   - For 12-core CPU: use `-workers 12`

2. Pattern length affects search speed:
   - Shorter patterns find matches faster
   - Prefix patterns (`123*`) are faster than contains patterns (`*123*`)
   - Complex patterns with multiple wildcards (`1*2*3*4`) take longer to match
   - Using multiple patterns increases search time proportionally
   - Special characters (`#` and `@`) are as fast as exact character matches

## Security Note

- Keep your private keys safe and secure
- Never share private keys of addresses containing real funds
- This tool is for educational and testing purposes only

## Technical Details

- Written in Go
- Uses ethereum-go libraries for key generation
- Thread-safe statistics tracking
- Advanced pattern matching algorithm with special character support
- Support for multiple simultaneous patterns
- Real-time performance monitoring

## GitHub Repository

Find the latest version and contribute to the project at:
https://github.com/grom42kem/genaddr

## Support the Project

If you find this tool useful, you can support its development by sending donations to:
`0x77777777b487e2FD60F3C60B080E03e7247338f6`

Thank you for your support! 