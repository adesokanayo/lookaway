import Image from "next/image";
import { latestDownloads } from "@/lib/releases";

export default async function Home() {
  const downloads = await latestDownloads();

  return (
    <main className="mx-auto flex min-h-full max-w-2xl flex-col px-6 py-24 sm:px-8 sm:py-32">
      <p className="text-[1.05rem] text-muted">Lookaway</p>
      <h1
        className="mt-8 text-[2.75rem] leading-[1.2] font-normal tracking-[-0.02em] sm:text-[3.35rem]"
        style={{ fontFamily: "var(--font-display), ui-serif, Georgia, serif" }}
      >
        Every 20 minutes, the screen takes over.
      </h1>
      <p className="mt-7 max-w-xl text-[1.25rem] leading-8 text-muted">
        Look about 20 feet away for 20 seconds. Notifications were too easy to swipe away. This
        actually takes the screen.
      </p>

      <div className="mt-12 flex flex-col gap-3 sm:flex-row">
        <a
          href={downloads.mac}
          className="inline-flex items-center justify-center rounded-md bg-accent px-7 text-[1.05rem] font-medium text-background"
          style={{ height: 52 }}
        >
          Download for Mac
        </a>
        <a
          href={downloads.windows}
          className="inline-flex items-center justify-center rounded-md border border-line px-7 text-[1.05rem] font-medium text-foreground"
          style={{ height: 52 }}
        >
          Download for Windows
        </a>
      </div>
      <p className="mt-5 text-[1rem] leading-6 text-muted">
        Free. Unsigned: Mac → Privacy &amp; Security → Open Anyway. Windows → More info → Run
        anyway.
      </p>

      <figure className="mt-24">
        <p className="mb-4 text-center text-[1.05rem] text-muted">
          This is the Small Eye icon that counts down the timer.
        </p>
        <Image
          src="/menu-bar.png"
          alt="Lookaway in the menu bar"
          width={360}
          height={80}
          className="mx-auto h-auto w-[220px]"
          priority
        />
      </figure>

      <figure className="mt-16">
        <p className="mb-4 text-center text-[1.05rem] text-muted">This is the menu.</p>
        <Image
          src="/menu.png"
          alt="Lookaway menu"
          width={720}
          height={480}
          className="mx-auto h-auto w-[340px]"
        />
      </figure>

      <figure className="mt-16">
        <p className="mb-4 text-center text-[1.05rem] text-muted">This is the Overlay you will see.</p>
        <Image
          src="/overlay.png"
          alt="Look away overlay"
          width={1440}
          height={900}
          className="mx-auto h-auto w-full rounded-sm"
        />
      </figure>

      <p className="mt-24 text-[1.2rem] leading-8 text-foreground/90">
        I built this because I sit on a laptop all day and my eyes get fried. The 20-20-20 rule is
        what the American Optometric Association and the American Academy of Ophthalmology
        recommend for digital eye strain. It eases symptoms. It is not a cure.
      </p>

      <footer className="mt-20 flex gap-6 text-[1rem] text-muted">
        <a href="https://github.com/adesokanayo/lookaway" className="hover:text-foreground">
          GitHub
        </a>
        <a
          href="https://github.com/adesokanayo/lookaway/releases/latest"
          className="hover:text-foreground"
        >
          Releases
        </a>
        <span>MIT</span>
      </footer>
    </main>
  );
}
