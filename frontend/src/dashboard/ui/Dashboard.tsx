import { Box, Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow } from "@material-ui/core";
import React from "react";

class Dashboard extends React.Component {
    public render() {
        return <Box m={2}>
            <TableContainer component={Paper}>
                <Table aria-label="simple table" size="small">
                    <TableHead>
                        <TableRow>
                            <TableCell>Job</TableCell>
                            <TableCell>Status</TableCell>
                            <TableCell>Git repository</TableCell>
                            <TableCell>Base</TableCell>
                            <TableCell>Pulls</TableCell>
                        </TableRow>
                    </TableHead>
                    <TableBody />
                </Table>
            </TableContainer>
        </Box>;
    }
}

export default Dashboard;