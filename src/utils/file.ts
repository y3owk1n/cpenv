import * as fs from "fs";
import * as util from "util";

export const copyFile = util.promisify(fs.copyFile);
