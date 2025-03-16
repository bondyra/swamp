import React, { memo } from 'react';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import { styled } from '@mui/material/styles';

import {JSONPath} from 'jsonpath-plus';


const NiceContainer = styled(TableContainer)(({ theme }) => ({
  color: theme.palette.success.main,
  '.MuiTableCell-root': {
    fontFamily: "monospace",
    fontSize: "10px",
    color: "white",
    margin: "0px",
    padding: "3px",
    border: "1px solid gray"
  },
}));

export default memo(({ data, selectedFields }) => {
  const multiple = Array.isArray(data)
  const fields = multiple ? ["metadata.id"].concat(selectedFields) : selectedFields;
  return (
    <NiceContainer>
      <Table size="small" aria-label="a dense table" sx={{mb: "5px"}}>
        { 
          multiple && 
          <TableHead>
            <TableRow key="header">
              {fields.map(f => <TableCell key={`header-${f}`}><b>{f.replace("obj.", "")}</b></TableCell>)}
            </TableRow>
          </TableHead>
        }
        <TableBody>
          {
            multiple &&
            data.map(d => 
                  <TableRow key={d.metadata.id}>
                    {fields.map(f => <TableCell key={`${d.metadata.id}-${f}`} align="left">{JSONPath({path: f, json: d})}</TableCell>)}
                </TableRow>
            )
          }
          {
            !multiple &&
            fields.map(f => 
                  <TableRow key={f}>
                    <TableCell key={`${f}-key`} align="left" sx={{fontSize: "10px"}}><b>{f}</b></TableCell>
                    <TableCell key={`${f}-val`} align="left">{JSONPath({path: f, json: data})}</TableCell>
                </TableRow>
            )
          }
        </TableBody>
      </Table>
    </NiceContainer>
  );
});
