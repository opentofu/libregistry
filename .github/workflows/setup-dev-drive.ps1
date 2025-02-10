# Configures a drive for testing in CI.

# When not using a GitHub Actions "larger runner", the `D:` drive is present and
# has similar or better performance characteristics than a ReFS dev drive.
# Sometimes using a larger runner is still more performant (e.g., when running
# the test suite) and we need to create a dev drive. This script automatically
# configures the appropriate drive.

# Note we use `Get-PSDrive` is not sufficient because the drive letter is assigned.
if (Test-Path "D:\") {
    Write-Output "Using existing drive at D:"
    $Drive = "D:"
    Write-Output `
	"DEV_DRIVE=$($Drive)" `
	>> $env:GITHUB_ENV
}

