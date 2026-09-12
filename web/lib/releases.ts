const repo = "adesokanayo/lookaway";
const latestPage = `https://github.com/${repo}/releases/latest`;

export type Downloads = {
  mac: string;
  windows: string;
};

type GithubAsset = {
  name: string;
  browser_download_url: string;
};

type GithubRelease = {
  assets?: GithubAsset[];
};

export async function latestDownloads(): Promise<Downloads> {
  try {
    const res = await fetch(`https://api.github.com/repos/${repo}/releases/latest`, {
      next: { revalidate: 3600 },
      headers: { Accept: "application/vnd.github+json" },
    });
    if (!res.ok) {
      return { mac: latestPage, windows: latestPage };
    }
    const body = (await res.json()) as GithubRelease;
    const assets = body.assets ?? [];
    const mac = assets.find((a) => a.name.endsWith(".dmg"))?.browser_download_url;
    const windows = assets.find((a) => a.name.endsWith(".zip") && a.name.includes("windows"))
      ?.browser_download_url;
    return {
      mac: mac ?? latestPage,
      windows: windows ?? latestPage,
    };
  } catch {
    return { mac: latestPage, windows: latestPage };
  }
}
