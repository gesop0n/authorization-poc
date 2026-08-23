import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/")({
  component: Home,
});

function Home() {
  return (
    <main className="flex flex-1 flex-col items-center justify-center gap-4 bg-zinc-50 dark:bg-black">
      <h1 className="text-3xl font-semibold tracking-tight text-black dark:text-zinc-50">
        authorization-poc
      </h1>
      <p className="text-zinc-600 dark:text-zinc-400">
        Edit src/routes/index.tsx to get started.
      </p>
    </main>
  );
}
