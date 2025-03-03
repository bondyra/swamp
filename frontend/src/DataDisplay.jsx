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
    color: "white",
    padding: "6px 4px",
    border: "1px solid white"
  },
}));

export default memo(({ data, selectedFields }) => {
  const multiple = Array.isArray(data)
  const fields = multiple ? ["resource_id"].concat(selectedFields) : selectedFields;
  return (
    <NiceContainer>
      <Table size="small" aria-label="a dense table">
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
                  <TableRow key={d.resource_id}>
                    {fields.map(f => <TableCell key={`${d.resource_id}-${f}`} align="left">{JSONPath({path: f, json: d})}</TableCell>)}
                </TableRow>
            )
          }
          {
            !multiple &&
            fields.map(f => 
                  <TableRow key={f}>
                    <TableCell align="left">{f}</TableCell>
                    <TableCell align="left">{JSONPath({path: f, json: data})}</TableCell>
                </TableRow>
            )
          }
        </TableBody>
      </Table>
    </NiceContainer>
  );
});
