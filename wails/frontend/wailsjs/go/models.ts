export namespace main {

	export class PrintResult {
	    success: boolean;
	    reason?: string;
	    handled: boolean;

	    static createFrom(source: any = {}) {
	        return new PrintResult(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.reason = source["reason"];
	        this.handled = source["handled"];
	    }
	}
	export class Printer {
	    name: string;
	    displayName: string;
	
	    static createFrom(source: any = {}) {
	        return new Printer(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.displayName = source["displayName"];
	    }
	}

}
