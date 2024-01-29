class Netfetch < Formula
  desc "CLI tool to scan for network policies in Kubernetes clusters/namespaces and provide a score based on the amount of untargeted workloads"
  homepage "https://github.com/deggja/netfetch"

  if OS.mac?
    url "https://github.com/deggja/netfetch/releases/download/3.2.6/netfetch_3.2.6_darwin_amd64.tar.gz"
    sha256 "7a5e2e20904507020bcc4fc69706456b2f887c3ddf8980339d9b07deafd4feb9"
  elsif OS.linux?
    url "https://github.com/deggja/netfetch/releases/download/3.2.6/netfetch_3.2.6_linux_amd64.tar.gz"
    sha256 "208738f0e508458454cabe2c823d1e1bba2c389655ba5d2474cce523e0ca33bf"
  end

  def install
    bin.install "netfetch"
  end

  test do
    system "#{bin}/netfetch", "--version"
  end
end
