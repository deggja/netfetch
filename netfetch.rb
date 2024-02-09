class Netfetch < Formula
  desc "CLI tool to scan for network policies in Kubernetes clusters/namespaces and provide a score based on the amount of untargeted workloads"
  homepage "https://github.com/deggja/netfetch"

  if OS.mac?
    url "https://github.com/deggja/netfetch/releases/download/3.2.8/netfetch_3.2.8_darwin_amd64.tar.gz"
    sha256 "71b17743763290f800d44647e5ed981d1b11066320591bc6c4ed7053d985e9ea"
  elsif OS.linux?
    url "https://github.com/deggja/netfetch/releases/download/3.2.8/netfetch_3.2.8_linux_amd64.tar.gz"
    sha256 "7461d611f0647cb7f756c93a9dc8f6ac32f9fc3b7c5875eee325601450ebb5c7"
  end

  def install
    bin.install "netfetch"
  end

  test do
    system "#{bin}/netfetch", "version"
  end
end
