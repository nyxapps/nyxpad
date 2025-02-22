export namespace main {
	
	export class ReadFileOutput {
	    content: string;
	    file: string;
	
	    static createFrom(source: any = {}) {
	        return new ReadFileOutput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.content = source["content"];
	        this.file = source["file"];
	    }
	}

}

