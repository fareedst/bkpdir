# [IMPL-FILE_STATISTICS] [ARCH-FILE_STATISTICS] [REQ-OUT_002] [REQ-OUTPUT_FORMATTING]

## Summary contract

Gather os.Stat metadata into FileStatInfo with human-readable size and file-type classification for list and template output formatting. Backward compatibility maintained with existing format strings; inc and full commands behave consistently.

INPUT: file path
OUTPUT: FileStatInfo or error
DATA: size thresholds KB/MB/GB/TB, os.FileMode bits

## GATHER_FILE_STAT_INFO

SPEC-ID: IMPL-FILE_STATISTICS::GATHER_FILE_STAT_INFO
STEP T001: STAT path AND populate FileStatInfo with size, mtime, and type
STEP T002: EMIT statistics compatible with legacy format strings for inc and full command output

- [IMPL-FILE_STATISTICS] [ARCH-FILE_STATISTICS] [REQ-OUT_002] [REQ-OUTPUT_FORMATTING] — How: stat path and populate FileStatInfo with name, size, human size, mtime, mode, and type.

PROCEDURE GATHER_FILE_STAT_INFO(path):
  info, err = OS.Stat(path)
  IF err: RETURN error
  RETURN FileStatInfo{Path, Name, Size, formatHumanSize(Size), MTime, Mode, getFileType(info)}

## FORMAT_HUMAN_SIZE

SPEC-ID: IMPL-FILE_STATISTICS::FORMAT_HUMAN_SIZE
STEP T001: SCALE bytes to TB/GB/MB/KB/B display string

- [IMPL-FILE_STATISTICS] [ARCH-FILE_STATISTICS] [REQ-OUT_002] [REQ-OUTPUT_FORMATTING] — How: scale bytes to TB/GB/MB/KB/B with one decimal for large units.

PROCEDURE FORMAT_HUMAN_SIZE(size):
  IF size >= TB: RETURN sprintf("%.1fTB", size/TB)
  ELIF size >= GB: RETURN sprintf("%.1fGB", size/GB)
  ELIF size >= MB: RETURN sprintf("%.1fMB", size/MB)
  ELIF size >= KB: RETURN sprintf("%.1fKB", size/KB)
  RETURN sprintf("%dB", size)

## GET_FILE_TYPE

SPEC-ID: IMPL-FILE_STATISTICS::GET_FILE_TYPE
STEP T001: CLASSIFY file type from os.FileMode bits

- [IMPL-FILE_STATISTICS] [ARCH-FILE_STATISTICS] [REQ-OUT_002] [REQ-OUTPUT_FORMATTING] — How: classify regular, directory, symlink, device, pipe, socket, or other from mode bits.

PROCEDURE GET_FILE_TYPE(info):
  IF regular: RETURN "regular"
  IF directory: RETURN "directory"
  IF symlink/device/pipe/socket bits: RETURN matching label
  RETURN "other"
