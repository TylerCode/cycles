# Cycles Packaging & Distribution Plan

**Purpose:** Ensure Cycles works correctly across all major Linux package formats and distribution methods

---

## Current Status Assessment

### ✅ Working
- **Snap:** Currently functional
- Source build from repository

### ❌ Broken/Not Working
- **Flatpak:** Never successfully configured
- **AppImage:** Believed to be broken

### ✓ To Verify
- Plain executable distribution
- Individual distro packages (.deb, .rpm)

---

## Sprint 10.5: Packaging & Distribution (v1.2.0)

**Duration:** 2-3 days
**Priority:** High (Critical for wide adoption)
**Should be completed as part of or immediately after Sprint 10**

---

## Task Breakdown

### 1. Fix AppImage Packaging

#### Current Issues to Investigate
- [ ] Check `.github/workflows/appimage.yml` for errors
- [ ] Verify AppRun script functionality
- [ ] Test on clean system
- [ ] Identify missing dependencies

#### Implementation Steps

**A. Audit Existing AppImage Setup**
```bash
# Review current workflow
cat .github/workflows/appimage.yml

# Check for AppDir structure
# Verify linuxdeploy usage
# Test manual AppImage build
```

**B. Fix AppImage Configuration**

Create/update `appimage/AppRun`:
```bash
#!/bin/bash
HERE="$(dirname "$(readlink -f "${0}")")"
export PATH="${HERE}/usr/bin:${PATH}"
export LD_LIBRARY_PATH="${HERE}/usr/lib:${LD_LIBRARY_PATH}"
export XDG_DATA_DIRS="${HERE}/usr/share:${XDG_DATA_DIRS}"

# Fyne requires these for proper rendering
export FYNE_SCALE=1.0
export GDK_BACKEND=x11

exec "${HERE}/usr/bin/cycles" "$@"
```

**C. AppImage Desktop File**

Create `us.tylerc.cycles.desktop` (verify it's correct):
```desktop
[Desktop Entry]
Type=Application
Name=Cycles
Comment=Desktop CPU Monitor
Exec=cycles
Icon=us.tylerc.cycles
Categories=System;Monitor;
Terminal=false
```

**D. Build Script**

Create `appimage/build-appimage.sh`:
```bash
#!/bin/bash
set -e

# Build the application
go build -o cycles -ldflags="-s -w"

# Create AppDir structure
mkdir -p AppDir/usr/bin
mkdir -p AppDir/usr/share/applications
mkdir -p AppDir/usr/share/icons/hicolor/256x256/apps
mkdir -p AppDir/usr/share/metainfo

# Copy files
cp cycles AppDir/usr/bin/
cp icon.png AppDir/usr/share/icons/hicolor/256x256/apps/us.tylerc.cycles.png
cp us.tylerc.cycles.desktop AppDir/usr/share/applications/
cp us.tylerc.cycles.appdata.xml AppDir/usr/share/metainfo/

# Download linuxdeploy if needed
if [ ! -f linuxdeploy-x86_64.AppImage ]; then
    wget https://github.com/linuxdeploy/linuxdeploy/releases/download/continuous/linuxdeploy-x86_64.AppImage
    chmod +x linuxdeploy-x86_64.AppImage
fi

# Create AppImage
./linuxdeploy-x86_64.AppImage \
    --appdir AppDir \
    --output appimage \
    --desktop-file=AppDir/usr/share/applications/us.tylerc.cycles.desktop \
    --icon-file=AppDir/usr/share/icons/hicolor/256x256/apps/us.tylerc.cycles.png

echo "AppImage created successfully!"
```

**E. Testing Checklist**
- [ ] AppImage runs on Ubuntu 20.04, 22.04, 24.04
- [ ] AppImage runs on Fedora (latest)
- [ ] AppImage runs on Arch Linux
- [ ] Icon displays correctly in file manager
- [ ] Desktop integration works (double-click launch)
- [ ] No dependency errors
- [ ] Graphics render correctly (Fyne UI)

---

### 2. Implement Flatpak Support

#### Why Flatpak Has Been Difficult
- Fyne applications need special consideration in Flatpak
- OpenGL/graphics permissions required
- File system access for `/proc` reading
- Proper SDK selection

#### Implementation Steps

**A. Create Flatpak Manifest**

Create `flatpak/us.tylerc.cycles.yml`:
```yaml
app-id: us.tylerc.cycles
runtime: org.freedesktop.Platform
runtime-version: '23.08'
sdk: org.freedesktop.Sdk
sdk-extensions:
  - org.freedesktop.Sdk.Extension.golang

command: cycles

finish-args:
  # X11 + XShm access
  - --share=ipc
  - --socket=x11

  # OpenGL/graphics access (required for Fyne)
  - --device=dri

  # File system access for /proc reading
  - --filesystem=/proc:ro
  - --filesystem=/sys:ro

  # Network access (for network monitoring)
  - --share=network

  # Access to host system info
  - --filesystem=host:ro

modules:
  - name: cycles
    buildsystem: simple
    build-options:
      append-path: /usr/lib/sdk/golang/bin
      env:
        GOBIN: /app/bin
        GOROOT: /usr/lib/sdk/golang
    build-commands:
      - . /usr/lib/sdk/golang/enable.sh && go build -o cycles -ldflags="-s -w"
      - install -Dm755 cycles /app/bin/cycles
      - install -Dm644 icon.png /app/share/icons/hicolor/256x256/apps/us.tylerc.cycles.png
      - install -Dm644 us.tylerc.cycles.desktop /app/share/applications/us.tylerc.cycles.desktop
      - install -Dm644 us.tylerc.cycles.appdata.xml /app/share/metainfo/us.tylerc.cycles.appdata.xml
    sources:
      - type: git
        url: https://github.com/TylerCode/cycles
        branch: main
```

**B. Flatpak Build Script**

Create `flatpak/build-flatpak.sh`:
```bash
#!/bin/bash
set -e

echo "Building Flatpak for Cycles..."

# Install flatpak-builder if needed
if ! command -v flatpak-builder &> /dev/null; then
    echo "Error: flatpak-builder not found. Install it first."
    echo "  Ubuntu/Debian: sudo apt install flatpak-builder"
    echo "  Fedora: sudo dnf install flatpak-builder"
    exit 1
fi

# Add Flathub if not already added
flatpak remote-add --if-not-exists flathub https://flathub.org/repo/flathub.flatpakrepo

# Install required SDK
flatpak install -y flathub org.freedesktop.Platform//23.08 org.freedesktop.Sdk//23.08
flatpak install -y flathub org.freedesktop.Sdk.Extension.golang

# Build the Flatpak
flatpak-builder --force-clean build-dir flatpak/us.tylerc.cycles.yml

# Create repository
flatpak-builder --repo=repo --force-clean build-dir flatpak/us.tylerc.cycles.yml

# Build bundle for distribution
flatpak build-bundle repo cycles.flatpak us.tylerc.cycles

echo "Flatpak bundle created: cycles.flatpak"
echo "To install: flatpak install cycles.flatpak"
echo "To run: flatpak run us.tylerc.cycles"
```

**C. Testing Checklist**
- [ ] Flatpak builds without errors
- [ ] Can read CPU info from /proc
- [ ] Graphics render correctly (Fyne with OpenGL)
- [ ] All tabs display data correctly
- [ ] Settings persist correctly
- [ ] No permission errors in logs
- [ ] Works on different distros (test in VMs)

---

### 3. Snap Verification & Enhancement

#### Current Status
✅ Working - but let's verify and document

#### Verification Steps
```bash
# Test current snap
sudo snap install tylercode-cycles

# Verify interfaces
snap connections tylercode-cycles

# Check for any confined issues
snap logs tylercode-cycles
```

#### Enhancement: Ensure All Permissions

Update `snap/snapcraft.yaml` to ensure all required plugs:
```yaml
plugs:
  hardware-observe:
  system-observe:
  network-observe:
  mount-observe:
  process-control:  # For process management features
  opengl:           # For Fyne rendering
  x11:              # For GUI
```

#### Testing Checklist
- [ ] Snap installs correctly
- [ ] All monitoring features work
- [ ] Process management works (when implemented)
- [ ] No confined interface errors
- [ ] Auto-updates work correctly

---

### 4. Traditional Package Formats

#### Debian/Ubuntu (.deb)

**Create `packaging/debian/` structure:**
```
packaging/debian/
├── DEBIAN/
│   ├── control
│   ├── postinst
│   └── prerm
├── usr/
│   ├── bin/
│   │   └── cycles
│   └── share/
│       ├── applications/
│       │   └── us.tylerc.cycles.desktop
│       ├── icons/
│       │   └── hicolor/256x256/apps/
│       │       └── us.tylerc.cycles.png
│       └── metainfo/
│           └── us.tylerc.cycles.appdata.xml
```

**Create `packaging/build-deb.sh`:**
```bash
#!/bin/bash
set -e

VERSION=$(grep 'Version:' config.go | cut -d'"' -f2)
ARCH=$(dpkg --print-architecture)
PACKAGE_NAME="tylercode-cycles_${VERSION}_${ARCH}"

# Build the binary
go build -o cycles -ldflags="-s -w"

# Create package structure
mkdir -p "${PACKAGE_NAME}/DEBIAN"
mkdir -p "${PACKAGE_NAME}/usr/bin"
mkdir -p "${PACKAGE_NAME}/usr/share/applications"
mkdir -p "${PACKAGE_NAME}/usr/share/icons/hicolor/256x256/apps"
mkdir -p "${PACKAGE_NAME}/usr/share/metainfo"

# Copy files
cp cycles "${PACKAGE_NAME}/usr/bin/"
cp us.tylerc.cycles.desktop "${PACKAGE_NAME}/usr/share/applications/"
cp icon.png "${PACKAGE_NAME}/usr/share/icons/hicolor/256x256/apps/us.tylerc.cycles.png"
cp us.tylerc.cycles.appdata.xml "${PACKAGE_NAME}/usr/share/metainfo/"

# Create control file
cat > "${PACKAGE_NAME}/DEBIAN/control" << EOF
Package: tylercode-cycles
Version: ${VERSION}
Section: utils
Priority: optional
Architecture: ${ARCH}
Maintainer: Tyler <your-email@example.com>
Description: Desktop CPU Monitor
 Cycles provides real-time CPU, memory, disk, and network monitoring
 with graphical visualization.
Depends: libc6, libgl1, libx11-6, libxcursor1, libxrandr2, libxinerama1, libxi6
EOF

# Build package
dpkg-deb --build "${PACKAGE_NAME}"

echo "Debian package created: ${PACKAGE_NAME}.deb"
```

#### Fedora/RHEL (.rpm)

**Create `packaging/cycles.spec`:**
```spec
Name:           cycles
Version:        0.4.1
Release:        1%{?dist}
Summary:        Desktop CPU Monitor

License:        MIT
URL:            https://github.com/TylerCode/cycles
Source0:        %{name}-%{version}.tar.gz

BuildRequires:  golang >= 1.20
BuildRequires:  mesa-libGL-devel
BuildRequires:  libXcursor-devel
BuildRequires:  libXrandr-devel
BuildRequires:  libXinerama-devel
BuildRequires:  libXi-devel

Requires:       mesa-libGL
Requires:       libXcursor
Requires:       libXrandr
Requires:       libXinerama
Requires:       libXi

%description
Cycles provides real-time CPU, memory, disk, and network monitoring
with graphical visualization similar to Windows Task Manager.

%prep
%setup -q

%build
go build -o cycles -ldflags="-s -w"

%install
mkdir -p %{buildroot}%{_bindir}
mkdir -p %{buildroot}%{_datadir}/applications
mkdir -p %{buildroot}%{_datadir}/icons/hicolor/256x256/apps
mkdir -p %{buildroot}%{_datadir}/metainfo

install -m 755 cycles %{buildroot}%{_bindir}/
install -m 644 us.tylerc.cycles.desktop %{buildroot}%{_datadir}/applications/
install -m 644 icon.png %{buildroot}%{_datadir}/icons/hicolor/256x256/apps/us.tylerc.cycles.png
install -m 644 us.tylerc.cycles.appdata.xml %{buildroot}%{_datadir}/metainfo/

%files
%license LICENSE
%{_bindir}/cycles
%{_datadir}/applications/us.tylerc.cycles.desktop
%{_datadir}/icons/hicolor/256x256/apps/us.tylerc.cycles.png
%{_datadir}/metainfo/us.tylerc.cycles.appdata.xml

%changelog
* Thu Nov 17 2025 Tyler <your-email@example.com> - 0.4.1-1
- Initial RPM package
```

**Create `packaging/build-rpm.sh`:**
```bash
#!/bin/bash
set -e

# Ensure rpmbuild directory structure exists
mkdir -p ~/rpmbuild/{BUILD,RPMS,SOURCES,SPECS,SRPMS}

# Copy spec file
cp packaging/cycles.spec ~/rpmbuild/SPECS/

# Create source tarball
VERSION=$(grep 'Version:' config.go | cut -d'"' -f2)
mkdir -p cycles-${VERSION}
cp -r *.go go.mod go.sum icon.png *.desktop *.xml LICENSE cycles-${VERSION}/
tar czf ~/rpmbuild/SOURCES/cycles-${VERSION}.tar.gz cycles-${VERSION}
rm -rf cycles-${VERSION}

# Build RPM
rpmbuild -ba ~/rpmbuild/SPECS/cycles.spec

echo "RPM package created in ~/rpmbuild/RPMS/"
```

---

### 5. AppData/Metainfo Validation

The `us.tylerc.cycles.appdata.xml` file is critical for all package formats.

#### Validate Current File
```bash
# Install validator
sudo apt install appstream-util  # Debian/Ubuntu
sudo dnf install appstream         # Fedora

# Validate
appstream-util validate-relax us.tylerc.cycles.appdata.xml
```

#### Update if Needed
Ensure the AppData file includes:
- Valid screenshots
- Release information
- Proper metadata tags
- Content rating
- Update contact info

---

### 6. CI/CD Integration

#### Update `.github/workflows/` for All Formats

**Create `.github/workflows/packages.yml`:**
```yaml
name: Build All Packages

on:
  push:
    tags:
      - 'v*'
  workflow_dispatch:

jobs:
  appimage:
    runs-on: ubuntu-20.04
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Install dependencies
        run: |
          sudo apt-get update
          sudo apt-get install -y libgl1-mesa-dev libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev libglfw3-dev libxxf86vm-dev
      - name: Build AppImage
        run: |
          chmod +x appimage/build-appimage.sh
          ./appimage/build-appimage.sh
      - name: Upload AppImage
        uses: actions/upload-artifact@v3
        with:
          name: AppImage
          path: '*.AppImage'

  flatpak:
    runs-on: ubuntu-latest
    container:
      image: bilelmoussaoui/flatpak-github-actions:freedesktop-23.08
      options: --privileged
    steps:
      - uses: actions/checkout@v3
      - uses: flatpak/flatpak-github-actions/flatpak-builder@v6
        with:
          bundle: cycles.flatpak
          manifest-path: flatpak/us.tylerc.cycles.yml
      - name: Upload Flatpak
        uses: actions/upload-artifact@v3
        with:
          name: Flatpak
          path: cycles.flatpak

  deb:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Install dependencies
        run: |
          sudo apt-get update
          sudo apt-get install -y libgl1-mesa-dev libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev
      - name: Build DEB
        run: |
          chmod +x packaging/build-deb.sh
          ./packaging/build-deb.sh
      - name: Upload DEB
        uses: actions/upload-artifact@v3
        with:
          name: DEB
          path: '*.deb'

  rpm:
    runs-on: ubuntu-latest
    container:
      image: fedora:latest
    steps:
      - uses: actions/checkout@v3
      - name: Install dependencies
        run: |
          dnf install -y golang rpm-build mesa-libGL-devel libXcursor-devel libXrandr-devel libXinerama-devel libXi-devel
      - name: Build RPM
        run: |
          chmod +x packaging/build-rpm.sh
          ./packaging/build-rpm.sh
      - name: Upload RPM
        uses: actions/upload-artifact@v3
        with:
          name: RPM
          path: '~/rpmbuild/RPMS/*/*.rpm'
```

---

## Testing Matrix

| Format | Ubuntu 20.04 | Ubuntu 22.04 | Ubuntu 24.04 | Fedora 39 | Arch Linux | Notes |
|--------|--------------|--------------|--------------|-----------|------------|-------|
| AppImage | ✓ | ✓ | ✓ | ✓ | ✓ | Should work everywhere |
| Flatpak | ✓ | ✓ | ✓ | ✓ | ✓ | Universal |
| Snap | ✓ | ✓ | ✓ | ✓ | ✓ | Already working |
| .deb | ✓ | ✓ | ✓ | ✗ | ✗ | Debian-based only |
| .rpm | ✗ | ✗ | ✗ | ✓ | ✗ | Red Hat-based only |
| AUR | ✗ | ✗ | ✗ | ✗ | ✓ | Arch User Repository |

---

## Documentation Updates

### README.md Section
Add comprehensive installation section:
```markdown
## Installation

### Linux

#### Flatpak (Recommended - All Distributions)
```bash
flatpak install cycles.flatpak
flatpak run us.tylerc.cycles
```

#### Snap (All Distributions)
```bash
sudo snap install tylercode-cycles
```

#### AppImage (All Distributions)
```bash
# Download from releases page
chmod +x Cycles-x86_64.AppImage
./Cycles-x86_64.AppImage
```

#### Ubuntu/Debian
```bash
sudo dpkg -i tylercode-cycles_0.4.1_amd64.deb
```

#### Fedora/RHEL
```bash
sudo dnf install tylercode-cycles-0.4.1-1.x86_64.rpm
```

#### Arch Linux (AUR)
```bash
yay -S cycles-bin
```
```

---

## Deliverables Checklist

### Required Files
- [ ] `appimage/build-appimage.sh`
- [ ] `appimage/AppRun`
- [ ] `flatpak/us.tylerc.cycles.yml`
- [ ] `flatpak/build-flatpak.sh`
- [ ] `packaging/build-deb.sh`
- [ ] `packaging/build-rpm.sh`
- [ ] `packaging/cycles.spec`
- [ ] `.github/workflows/packages.yml`
- [ ] Updated `README.md` with installation instructions
- [ ] Updated `CHANGELOG.md` with package format additions

### Testing Complete
- [ ] AppImage works on 3+ distros
- [ ] Flatpak builds and runs without errors
- [ ] Snap still works (regression test)
- [ ] .deb installs on Ubuntu/Debian
- [ ] .rpm installs on Fedora
- [ ] All packages show icon correctly
- [ ] All packages integrate with desktop environment
- [ ] All packages can read system info correctly
- [ ] No permission/sandbox issues

### CI/CD
- [ ] GitHub Actions build all packages automatically
- [ ] Release workflow uploads all package types
- [ ] Package naming is consistent
- [ ] Version numbers sync across all formats

---

## Success Criteria

1. **AppImage:** Works out-of-the-box on major distributions without dependency issues
2. **Flatpak:** Builds cleanly and runs with proper permissions for system monitoring
3. **Snap:** Continues to work as it currently does (no regression)
4. **Traditional Packages:** .deb and .rpm install cleanly on their respective platforms
5. **Documentation:** Clear installation instructions for each package format
6. **Automation:** CI/CD builds all package formats on release

---

## Common Issues & Solutions

### Flatpak Issues

**Problem:** Can't read /proc filesystem
**Solution:** Add `--filesystem=/proc:ro` to finish-args

**Problem:** Graphics not rendering (black screen)
**Solution:** Ensure `--device=dri` and `--socket=x11` permissions

**Problem:** Build fails with Go errors
**Solution:** Verify Go SDK extension is installed and enabled

### AppImage Issues

**Problem:** "Cannot execute binary file"
**Solution:** Ensure AppRun has execute permissions and proper shebang

**Problem:** Missing libraries
**Solution:** Use linuxdeploy with --bundle-non-qt-libs flag

**Problem:** Icon not showing
**Solution:** Verify .desktop file and icon paths are correct

### General Issues

**Problem:** Version number mismatch
**Solution:** Create single source of truth (config.go) and script version extraction

**Problem:** Different behavior in packaged vs source
**Solution:** Test with relative paths, ensure icon.png is embedded or properly located

---

## Timeline Integration

This packaging work should be integrated into **Sprint 10** or done as **Sprint 10.5** before the v1.2.0 release.

**Recommended Approach:**
- Fix AppImage: Day 1
- Implement Flatpak: Day 2
- Create traditional packages: Day 3
- CI/CD integration: Day 4
- Testing on all platforms: Day 5

---

*This plan ensures Cycles can be easily installed by users on any Linux distribution using their preferred package format.*
