import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";

import { EnhancedNodesHeader } from "@/components/nodes/EnhancedNodesHeader";

describe("EnhancedNodesHeader", () => {
  it("defaults to Silmari copy in the subtitle", () => {
    render(
      <EnhancedNodesHeader
        totalNodes={3}
        onlineCount={2}
        offlineCount={1}
        degradedCount={0}
        startingCount={0}
        isConnected={true}
      />,
    );

    expect(screen.getByText("Agent Nodes")).toBeInTheDocument();
    expect(
      screen.getByText(
        "Monitor and manage your AI agent nodes in the Silmari control plane.",
      ),
    ).toBeInTheDocument();
  });
});
