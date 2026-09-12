import Image from "next/image";
import { latestDownloads } from "@/lib/releases";

export default async function Home() {
  const downloads = await latestDownloads();

  return (
    <main className="mx-auto flex min-h-full max-w-xl flex-col px-6 py-20 sm:py-28">
      <p className="text-[13px] tracking-[0.18em] text-muted uppercase">Lookaway</p>
      <h1
        className="mt-6 text-[2.15rem] leading-[1.15] font-medium tracking-[-0.03em] sm:text-[2.6rem]"
        style={{ fontFamily: "var(--font-display), ui-serif, Georgia, serif" }}
      >
        Every 20 minutes, the screen takes over.
      </h1>
      <p className="mt-5 max-w-md text-[1.05rem] leading-7 text-muted">
        Look about 20 feet away for 20 seconds. Notifications were too easy to swipe away. This
        actually takes the screen.
      </p>

      <div className="mt-10 flex flex-col gap-3 sm:flex-row">
        <a
          href={downloads.mac}
          className="inline-flex h-11 items-center justify-center rounded-full bg-accent px-6 text-[15px] font-medium text-background"
        >
          Download for Mac
        </a>
        <a
          href={downloads.windows}
          className="inline-flex h-11 items-center justify-center rounded-full border border-line px-6 text-[15px] font-medium text-foreground"
        >
          Download for Windows
        </a>
      </div>
      <p className="mt-4 text-[13px] leading-5 text-muted">
        Free. Unsigned: Mac → Privacy &amp; Security → Open Anyway. Windows → More info → Run
        anyway.
      </p>

      <figure className="mt-20">
        <p className="mb-3 text-center text-[13px] text-muted">
          This is the Small Eye icon that counts down the timer.
        </p>
        <Image
          src="/menu-bar.png"
          alt="Lookaway in the menu bar"
          width={360}
          height={80}
          className="mx-auto h-auto w-[180px]"
          priority
        />
      </figure>

      <figure className="mt-14">
        <p className="mb-3 text-center text-[13px] text-muted">This is the menu.</p>
        <Image
          src="/menu.png"
          alt="Lookaway menu"
          width={720}
          height={480}
          className="mx-auto h-auto w-[280px]"
        />
      </figure>

      <figure className="mt-14">
        <p className="mb-3 text-center text-[13px] text-muted">This is the Overlay you will see.</p>
        <Image
          src="/overlay.png"
          alt="Look away overlay"
          width={1440}
          height={900}
          className="mx-auto h-auto w-full rounded-sm"
        />
      </figure>

      <p className="mt-20 text-[1.02rem] leading-7 text-foreground/90">
        I built this because I sit on a laptop all day and my eyes get fried. The 20-20-20 rule is
        what the American Optometric Association and the American Academy of Ophthalmology
        recommend for digital eye strain. It eases symptoms. It is not a cure.
      </p>

      <footer className="mt-16 flex gap-5 text-[13px] text-muted">
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
