

import React, { memo } from 'react';
import SingleFieldPicker from './pickers/SingleFieldPicker';


export default memo(({ op, change }) => {
    const options = [
        {value: "==", description: "Left equals to right"},
        {value: "!=", description: "Left not equals to right"},
        {value: "~~", description: "Left like right"},
        {value: "!~", description: "Left not like right"},
    ]
  return (
    <SingleFieldPicker value={op} updateData={(v) => {change(v)}} options={options}
    valuePlaceholder="op" popperPrompt="Select operation"/>
  );
});
