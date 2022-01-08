import { Paper } from "@mui/material";
import React from "react";

export function TestBreakPoints(): JSX.Element {
	return (
		<Paper
			sx={{
				fontSize: {
					xs: "14px",
					sm: "24px",
					md: "36px",
					lg: "48px",
					xl: "72px",
				},
				width: {
					xs: "100px",
					sm: "160px",
					md: "240px",
					lg: "360px",
					xl: "480px",
				},
			}}
		>
			paper
		</Paper>
	);
}
