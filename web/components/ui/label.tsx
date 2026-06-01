import type * as React from "react";
import { cn } from "@/lib/utils";

function Label({
  className,
  htmlFor,
  ...props
}: React.ComponentProps<"label"> & { htmlFor: string }) {
  return (
    // biome-ignore lint/a11y/noLabelWithoutControl: Shared label component receives htmlFor from callers.
    <label
      data-slot="label"
      htmlFor={htmlFor}
      className={cn(
        "flex items-center gap-2 text-sm font-medium leading-none text-foreground",
        className,
      )}
      {...props}
    />
  );
}

export { Label };
