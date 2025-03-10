# Ethereum Address Generator (genaddr.exe)

A fast and flexible Ethereum address generator with pattern matching capabilities. Generate Ethereum addresses that match specific patterns at the beginning, end, or anywhere in the address.

## Features

- Pattern-based address generation
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
  - Pattern to match. Examples:
    - `123*` : address starts with "123"
    - `*123` : address ends with "123"
    - `*123*` : address contains "123"
- `-workers int` (default: 4)
  - Number of worker goroutines
- `-continue`
  - Continue searching after finding a match
- `-output string`
  - Save found addresses to file
- `-help`
  - Show help message

### Examples

1. Find an address starting with "123":
```bash
genaddr.exe -pattern "123*"
```

2. Find multiple addresses containing "dead" and save them:
```bash
genaddr.exe -pattern "*dead*" -continue -output results.txt
```

3. Use 8 worker threads to find an address ending with "cafe":
```bash
genaddr.exe -pattern "*cafe" -workers 8
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

## Security Note

- Keep your private keys safe and secure
- Never share private keys of addresses containing real funds
- This tool is for educational and testing purposes only

## Technical Details

- Written in Go
- Uses ethereum-go libraries for key generation
- Thread-safe statistics tracking
- Efficient pattern matching algorithms
- Real-time performance monitoring 