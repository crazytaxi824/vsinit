import React from "react";
import { styled } from "@mui/material/styles";
import { Paper, Box } from "@mui/material";

// styled 标准元素
const StyledDiv = styled("div")(({ theme }) => ({
	margin: theme.spacing(2),
	backgroundColor: theme.palette.success.main,
	color: theme.palette.success.contrastText,
	width: "100px",
	height: "80px",
	lineHeight: "80px",
	textAlign: "center",
}));

const StyledInput = styled("input")(({ theme }) => ({
	margin: theme.spacing(2),
	backgroundColor: theme.palette.secondary.main,
	color: theme.palette.secondary.contrastText,
	width: "100px",
	height: "30px",
	lineHeight: "30px",
	textAlign: "center",
}));

// styled mui 组件
const StyledPaper = styled(Paper)(({ theme }) => ({
	margin: theme.spacing(2),
	backgroundColor: theme.palette.background.paper,
	color: theme.palette.text.primary,
	width: "100px",
	height: "60px",
	lineHeight: "60px",
	textAlign: "center",
}));

const StyledFlexBox = styled(Box)(({ theme }) => ({
	margin: theme.spacing(2),
	backgroundColor: theme.palette.error.main,
	color: theme.palette.error.contrastText,
	display: "flex",
	justifyContent: "center",
	alignItems: "center",
}));

const StyledPaper2 = styled(Paper)((props) => {
	console.log(props.className);

	return {
		width: "100px",
		height: "30px",
		lineHeight: "30px",
		textAlign: "center",
	};
});

export function MyStyledComponent(): JSX.Element {
	return (
		<>
			<StyledDiv>div</StyledDiv>
			<StyledInput />
			<StyledPaper>paper</StyledPaper>
			<StyledPaper2 className="foo">paper 2</StyledPaper2>
			<StyledFlexBox>
				<StyledDiv>div</StyledDiv>
				<StyledDiv>div</StyledDiv>
				<StyledDiv>div</StyledDiv>
			</StyledFlexBox>
		</>
	);
}
