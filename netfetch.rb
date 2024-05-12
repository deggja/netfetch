class Netfetch < Formula
  desc "CLI tool to scan for network policies in Kubernetes clusters/namespaces and provide a score based on the amount of untargeted workloads"
  homepage "https://github.com/deggja/netfetch"

  if OS.mac?
    url "https://github.com/deggja/netfetch/releases/download/5.2.2/netfetch_5.2.2_darwin_amd64.tar.gz"
    sha256 "4a8687706fb9ca8df94f020c27540c3296c3032a5e0f6c0e1427b24438d8d3e7"
  elsif OS.linux?
    url "https://github.com/deggja/netfetch/releases/download/5.2.2/netfetch_5.2.2_linux_amd64.tar.gz"
    sha256 "7f220340524cf8bb8822ac589bdd8ad83e5e27a65cf8083b06b4cce78220f212"
  end

  def install
    bin.install "netfetch"
  end

  test do
    system "#{bin}/netfetch", "version"
  end
end
