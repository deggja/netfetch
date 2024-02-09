class Netfetch < Formula
  desc "CLI tool to scan for network policies in Kubernetes clusters/namespaces and provide a score based on the amount of untargeted workloads"
  homepage "https://github.com/deggja/netfetch"

  if OS.mac?
    url "https://github.com/deggja/netfetch/releases/download/3.2.11/netfetch_3.2.11_darwin_amd64.tar.gz"
    sha256 "2f9954637f21af44d7d7c1ed9b503195a38d19981b2ef8657ae01a6d7928d8d7"
  elsif OS.linux?
    url "https://github.com/deggja/netfetch/releases/download/3.2.11/netfetch_3.2.11_linux_amd64.tar.gz"
    sha256 "985d48deaaccbcd26c28393734533097e0130e354fdb5a02529bde48e2787e49"
  end

  def install
    bin.install "netfetch"
  end

  test do
    system "#{bin}/netfetch", "version"
  end
end
