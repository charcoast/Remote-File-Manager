import { Grid, IconButton, styled } from "@mui/material"

const CircleButton = styled(IconButton)({
    width: "1rem",
    height: "1rem",
    borderRadius: "50%",
    backgroundColor: "#d1d5db",
    "&:hover": {
        backgroundColor: "#e0e0e0",
    },
});

export const MacButtons = () => {
    return (
        <Grid container spacing={1} sx={{ marginBottom: "5%" }}>
            <Grid item>
                <CircleButton
                    sx={{
                        backgroundColor: "#FF464F",
                        "&:hover": { backgroundColor: "#FF464F" },
                    }}
                />
            </Grid>
            <Grid item>
                <CircleButton
                    sx={{
                        backgroundColor: "#FFB41B",
                        "&:hover": { backgroundColor: "#FFB41B" },
                    }}
                />
            </Grid>
            <Grid item>
                <CircleButton
                    sx={{
                        backgroundColor: "#20CA2D",
                        "&:hover": { backgroundColor: "#20CA2D" },
                    }}
                />
            </Grid>
        </Grid>
    )
}