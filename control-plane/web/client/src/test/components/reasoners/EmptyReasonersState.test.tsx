import { render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";

import { EmptyReasonersState } from "@/components/reasoners/EmptyReasonersState";

describe("EmptyReasonersState", () => {
  it("uses Silmari branding in the no-reasoners getting-started copy", () => {
    render(
      <EmptyReasonersState
        type="no-reasoners"
        onRefresh={vi.fn()}
      />,
    );

    expect(screen.getByText("No Reasoners Available")).toBeInTheDocument();
    expect(
      screen.getByText(
        "Launch an agent node to register reasoners with Silmari. They will appear here as soon as they are online.",
      ),
    ).toBeInTheDocument();
  });
});
