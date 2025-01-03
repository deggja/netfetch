class Netfetch < Formula
  desc "CLI tool to scan for network policies in Kubernetes clusters/namespaces and provide a score based on the amount of untargeted workloads"
  homepage "https://github.com/deggja/netfetch"

  if OS.mac?
    url "https://github.com/deggja/netfetch/releases/download/v0.5.4/netfetch_0.5.4_darwin_amd64.tar.gz"
    sha256 "ada24b740c746bdf14e67c2153dbc02440462aa765e7ac31b87439aff845d48c"
  elsif OS.linux?
    url "https://github.com/deggja/netfetch/releases/download/v0.5.4/netfetch_0.5.4_linux_amd64.tar.gz"
    sha256 "8757efca2f1196777acc45299773da105d8ce40a260e5de8d9f72d942f6f896b"
  end

  def install
    bin.install "netfetch"
  end

  test do
    system "#{bin}/netfetch", "version"
  end
end
