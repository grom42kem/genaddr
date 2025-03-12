# Ethereum Address Generator (genaddr.exe)

A fast and flexible Ethereum address generator with pattern matching capabilities. Generate Ethereum addresses that match specific patterns at the beginning, end, or anywhere in the address.

## Features

- Advanced pattern-based address generation
- Support for multiple patterns
- Complex pattern matching with multiple wildcards
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
  - Multiple patterns:
    - `123*,*456,*789*` : matches any of these patterns
    - `dead*,*cafe*,babe*` : matches addresses starting with "dead" OR containing "cafe" OR starting with "babe"
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

3. Find an address with sequential numbers:
```bash
genaddr.exe -pattern "1*2*3*4*5" -workers 4
```

4. Complex pattern matching:
```bash
genaddr.exe -pattern "aa*bb*cc" -workers 8
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

## Security Note

- Keep your private keys safe and secure
- Never share private keys of addresses containing real funds
- This tool is for educational and testing purposes only

## Technical Details

- Written in Go
- Uses ethereum-go libraries for key generation
- Thread-safe statistics tracking
- Advanced pattern matching algorithm
- Support for multiple simultaneous patterns
- Real-time performance monitoring 