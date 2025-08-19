# Your Steampipe install root (default)
INSTALL_ROOT="${STEAMPIPE_INSTALL_DIR:-$HOME/.steampipe}"

# Create the local plugin dir (no hub.steampipe.io, no version dir)
mkdir -p "$INSTALL_ROOT/plugins/local/aws"

# Move/rename so the file is named aws.plugin and is the ONLY *.plugin here
mv "$INSTALL_ROOT/plugins/hub.steampipe.io/local/aws/dev/aws.plugin" \
  "$INSTALL_ROOT/plugins/local/aws/aws.plugin"

# Make sure it's executable
chmod +x "$INSTALL_ROOT/plugins/local/aws/aws.plugin"
