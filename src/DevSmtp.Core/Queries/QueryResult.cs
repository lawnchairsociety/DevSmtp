namespace DevSmtp.Core.Queries
{
    public abstract class QueryResult
    {
        public QueryResult()
        {
            this.Succeeded = true;
        }

        public QueryResult(Exception error)
        {
            this.Error = error;
            this.Succeeded = false;
        }

        public bool Succeeded { get; }
        public Exception? Error { get; }
    }
}
