using DevSmtp.Core.Stores;

namespace DevSmtp.Core.Commands
{
    public sealed class MailHandler : ICommandHandler<Mail, MailResult>
    {
        private readonly IDataStore _dataStore;

        public MailHandler(IDataStore dataStore)
        {
            this._dataStore = dataStore ?? throw new ArgumentNullException(nameof(dataStore));
        }

        public Task<MailResult> ExecuteAsync(Mail command, CancellationToken cancellationToken = default)
        {
            throw new NotImplementedException();
        }
    }
}
