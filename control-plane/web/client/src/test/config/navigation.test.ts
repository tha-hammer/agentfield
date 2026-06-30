import { describe, expect, it } from "vitest";

import { resourceLinks } from "@/config/navigation";

describe("navigation resource links", () => {
  it("keeps legacy link targets while rebranding visible labels", () => {
    expect(resourceLinks).toEqual([
      expect.objectContaining({
        title: "Silmari Docs",
        href: "https://agentfield.ai/docs",
      }),
      expect.objectContaining({
        title: "Silmari GitHub",
        href: "https://github.com/Agent-Field/agentfield/",
      }),
    ]);
  });
});
