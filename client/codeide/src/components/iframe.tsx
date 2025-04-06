import React from "react";
import { AppBar, Box, Button, Toolbar } from "@mui/material";
import { styled } from "@mui/system";
import { remoteserverurl } from "./../libs/Url";

const IframeContainer = styled(Box)(({ theme }) => ({
  width: "100vw",
  height: "500px",
  boxShadow: "0 4px 20px rgba(0, 0, 0, 0.2)",
  borderRadius: "10px",
  overflow: "hidden",
  backgroundColor: "white",
  display: "flex",
  flexDirection: "row",
  margin: "0 auto", // Center the container
  marginTop: "50px", // Add some top margin
}));

const Header = styled(AppBar)(({ theme }) => ({
  backgroundColor: "#f8f8f8",
  boxShadow: "none",
  borderBottom: "1px solid #ddd",
}));

const DotsContainer = styled(Box)({
  display: "flex",
  marginRight: "auto",
});

const Dot = styled(Box)(({ color }) => ({
  width: "12px",
  height: "12px",
  borderRadius: "50%",
  marginRight: "5px",
  backgroundColor: color,
}));

const Iframe = () => {
  const handleGoBack = () => {
  };

  return (
    <IframeContainer>
      <Header position="static">
        <Toolbar>
          <DotsContainer>
            <Dot color="yellow" /> {/* Added yellow dot */}
            <Dot color="red" />
            <Dot color="green" />
          </DotsContainer>
          <Button variant="contained" color="primary" onClick={handleGoBack}>
            X
          </Button>
        </Toolbar>
      </Header>
      <iframe
        src={"remoteserverurl"}
        style={{ flex: 1, width: "100%", border: "none" }}
      />
    </IframeContainer>
  );
};

export default Iframe;
